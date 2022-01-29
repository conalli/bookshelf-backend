package db

import (
	"context"
	"log"
	"os"

	"github.com/conalli/bookshelf-backend/models/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	if os.Getenv("LOCAL") == "dev" {
		db = os.Getenv("DEV_DB_NAME")
	} else if os.Getenv("LOCAL") == "test" {
		db = os.Getenv("TEST_DB_NAME")
	} else {
		db = os.Getenv("DB_NAME")
	}
	collection := c.DB.Database(db).Collection(collectionName)
	return *collection
}

// SessionWithTransaction takes a context and transaction func and returns the result of the transaction.
func (c *Client) SessionWithTransaction(ctx context.Context, transactionFunc func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := c.DB.StartSession(opts)
	defer sess.EndSession(ctx)
	if err != nil {
		log.Println("could not start db session")
		return nil, errors.NewInternalServerError()
	}
	txnOpts := options.Transaction().SetReadPreference(readpref.Primary())
	res, err := sess.WithTransaction(ctx, transactionFunc, txnOpts)
	return res, err
}
