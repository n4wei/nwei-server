package mongo

import (
	"context"
	"strings"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/n4wei/nwei-server/lib/logger"
)

const (
	DBClientContextKey = "dbclient"
)

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

func (c *DBClient) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *DBClient) Create(ctx context.Context, collectionName string, data interface{}) error {
	result, err := c.client.Database(c.dbName).Collection(collectionName).InsertOne(ctx, data)
	if err != nil {
		return err
	}

	c.logger.Printf("inserted %v into %s.%s with result %v", data, c.dbName, collectionName, *result)
	return nil
}

func (c *DBClient) List(ctx context.Context, collectionName string, result interface{}, handleResult func(interface{}) error) error {
	cursor, err := c.client.Database(c.dbName).Collection(collectionName).Find(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			c.logger.Error(err)
		}
	}()

	for cursor.Next(ctx) {
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
