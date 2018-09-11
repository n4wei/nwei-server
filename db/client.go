package db

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/n4wei/nwei-server/lib/logger"
)

type Client interface {
	Create(string, string, interface{}) error
	List(string, string, interface{}, func(interface{}) error) error
	Close() error
}

type DBConfig struct {
	URL  string
	Port int

	Username string
	Password string
	AuthDB   string

	Logger logger.Logger
}

type DBClient struct {
	client *mongo.Client
	logger logger.Logger
}

func NewClient(config DBConfig) (*DBClient, error) {
	auth := ""
	if config.Username != "" && config.Password != "" {
		auth = fmt.Sprintf("%s:%s@", config.Username, config.Password)
	}
	authDB := ""
	if config.AuthDB != "" {
		authDB = fmt.Sprintf("/%s", config.AuthDB)
	}
	addr := fmt.Sprintf("mongodb://%s%s:%d%s", auth, config.URL, config.Port, authDB)

	client, err := mongo.Connect(nil, addr)
	if err != nil {
		return nil, err
	}

	config.Logger.Printf("db listening on :%d", config.Port)

	return &DBClient{
		client: client,
		logger: config.Logger,
	}, nil
}

func (c *DBClient) Close() error {
	return c.client.Disconnect(nil)
}

// TODO: timeouts
func (c *DBClient) Create(dbName string, collectionName string, data interface{}) error {
	result, err := c.client.Database(dbName).Collection(collectionName).InsertOne(nil, data)
	if err != nil {
		return err
	}

	c.logger.Printf("inserted %v into %s.%s with result %v", data, dbName, collectionName, *result)
	return nil
}

// TODO: timeouts
func (c *DBClient) List(dbName string, collectionName string, result interface{}, handleResult func(interface{}) error) error {
	cursor, err := c.client.Database(dbName).Collection(collectionName).Find(nil, nil)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(nil) {
		err := cursor.Decode(result)
		if err != nil {
			return err
		}

		err = handleResult(result)
		if err != nil {
			return err
		}
	}

	return cursor.Err()
}
