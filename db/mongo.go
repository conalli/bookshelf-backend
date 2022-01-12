package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client represents a mongo db client.
type Client struct {
	DB *mongo.Client
}

// NewMongoClient uses a context to create a new client connection based on the MONGO_URI env var.
func NewMongoClient(ctx context.Context) *Client {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		DB: client,
	}
}

// MongoCollection uses the DB_NAME env var, and returns a collection based on the collectionName and client.
func (c *Client) MongoCollection(collectionName string) mongo.Collection {
	var db string
	if os.Getenv("LOCAL") == "true" {
		db = os.Getenv("DEV_DB_NAME")
	} else {
		db = os.Getenv("DB_NAME")
	}
	collection := c.DB.Database(db).Collection(collectionName)
	return *collection
}
