package db

import (
	"context"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/n4wei/nwei-server/lib/logger"
)

const (
	defaultTimeout = 3 * time.Second
)

type Client interface {
	Create(string, interface{}) error
	List(string, interface{}, func(interface{}) error) error
	Close() error
}

type DBConfig struct {
	URL    string
	Logger logger.Logger
}

type DBClient struct {
	client *mongo.Client
	dbName string
	logger logger.Logger
}

func NewClient(config DBConfig) (*DBClient, error) {
	client, err := mongo.Connect(context.Background(), config.URL)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(config.URL, "/")
	return &DBClient{
		client: client,
		dbName: parts[len(parts)-1],
		logger: config.Logger,
	}, nil
}

func (c *DBClient) Close() error {
	return c.client.Disconnect(context.Background())
}

func (c *DBClient) Create(collectionName string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	result, err := c.client.Database(c.dbName).Collection(collectionName).InsertOne(ctx, data)
	if err != nil {
		return err
	}

	c.logger.Printf("inserted %v into %s.%s with result %v", data, c.dbName, collectionName, *result)
	return nil
}

func (c *DBClient) List(collectionName string, result interface{}, handleResult func(interface{}) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	cursor, err := c.client.Database(c.dbName).Collection(collectionName).Find(ctx, nil)
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
