package repositories

import (
	"context"
	"github.com/inkoba/app_for_HR/internal/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	client             *mongo.Client
	collection         *mongo.Collection
	salariesCollection *mongo.Collection
	logger             *logrus.Logger
}

func NewMongoConfig(c config.Config, logger *logrus.Logger) *MongoConfig {
	// Create client
	clientOptions := options.Client().ApplyURI(c.URI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Fatal(err)
	}

	collection := client.Database(c.Database).Collection("users")
	salariesCollection := client.Database(c.Database).Collection("salaries")

	return &MongoConfig{client, collection, salariesCollection, logger}
}

func (c MongoConfig) Ping() error {
	err := c.client.Ping(context.TODO(), nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	return err
}
