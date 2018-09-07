package weight

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"GET": "foo"}`))
	case "POST":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"POST": "bar"}`))
	default:
		http.Error(w, "only GET and POST are allowed", http.StatusBadRequest)
	}
}
