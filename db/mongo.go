package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient uses a context to create a new client connection based on the MONGO_URI env var.
func MongoClient(ctx context.Context) mongo.Client {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return *client
}

// MongoCollection uses the DB_NAME env var, and returns a collection based on the collectionName and client.
func MongoCollection(client mongo.Client, collectionName string) mongo.Collection {
	db := os.Getenv("DB_NAME")
	collection := client.Database(db).Collection(collectionName)
	return *collection
}
