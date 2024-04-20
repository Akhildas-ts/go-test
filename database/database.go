package database

import (
	"context"
	"fmt"
	"lock/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase(cfg config.Config) (*mongo.Database, error) {
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the MongoDB server
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB server: %v", err)
	}

	// make these as a monog db 
	// Set global variable for MongoDB database
	DB = client.Database(cfg.MongoDBName)

	return DB, nil
}
