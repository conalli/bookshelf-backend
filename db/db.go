package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoContext() (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancelFunc
}

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

func MongoCollection(client mongo.Client, collectionName string) mongo.Collection {
	db := os.Getenv("DB_NAME")
	collection := client.Database(db).Collection(collectionName)
	return *collection
}
