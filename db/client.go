package db

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/n4wei/nwei-server/lib/logger"
)

const (
	defaultTimeout = 3 * time.Second
)

type Client interface {
	Create(string, string, interface{}) error
	List(string, string, interface{}, func(interface{}) error) error
	Close() error
}

type DBConfig struct {
	URL    string
	Logger logger.Logger
}

type DBClient struct {
	client *mongo.Client
	logger logger.Logger
}

func NewClient(config DBConfig) (*DBClient, error) {
	client, err := mongo.Connect(context.Background(), config.URL)
	if err != nil {
		return nil, err
	}

	return &DBClient{
		client: client,
		logger: config.Logger,
	}, nil
}

func (c *DBClient) Close() error {
	return c.client.Disconnect(context.Background())
}

func (c *DBClient) Create(dbName string, collectionName string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	result, err := c.client.Database(dbName).Collection(collectionName).InsertOne(ctx, data)
	if err != nil {
		return err
	}

	c.logger.Printf("inserted %v into %s.%s with result %v", data, dbName, collectionName, *result)
	return nil
}

func (c *DBClient) List(dbName string, collectionName string, result interface{}, handleResult func(interface{}) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	cursor, err := c.client.Database(dbName).Collection(collectionName).Find(ctx, nil)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
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
