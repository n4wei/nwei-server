package weight

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/n4wei/nwei-server/db"
	"github.com/n4wei/nwei-server/db/mongo"
)

const (
	dbCollection = "weight"
	timeout      = 5 * time.Second
)

type Weight struct {
	Value float64 `json:"value" bson:"value"`
	Time  int64   `json:"time" bson:"time"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	default:
		http.Error(w, "only GET and POST are allowed", http.StatusBadRequest)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	var weight Weight
	var weights []Weight

	dbClient, ok := r.Context().Value(mongo.DBClientContextKey).(db.Client)
	if !ok {
		handleErr(w, errors.New("no DB client available"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := dbClient.List(ctx, dbCollection, &weight, func(result interface{}) error {
		r, ok := result.(*Weight)
		if !ok {
			return errors.New("could not convert database entry to type Weight")
		}
		weights = append(weights, *r)
		return nil
	})
	if err != nil {
		handleErr(w, err)
	}

	data, err := json.Marshal(weights)
	if err != nil {
		handleErr(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleErr(w, err)
	}
	defer r.Body.Close()

	var weight Weight
	err = json.Unmarshal(data, &weight)
	if err != nil {
		handleErr(w, err)
	}
	weight.Time = time.Now().Unix()

	dbClient, ok := r.Context().Value(mongo.DBClientContextKey).(db.Client)
	if !ok {
		handleErr(w, errors.New("no DB client available"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = dbClient.Create(ctx, dbCollection, weight)
	if err != nil {
		handleErr(w, err)
	}

	w.WriteHeader(http.StatusCreated)
}

func handleErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
