package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DelAccReq represents the request body for the DelAcc endpoint.
type DelAccReq struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
}

// DelAccRes represents the response from the DelAcc endpoint.
type DelAccRes struct {
	NumDeleted int    `json:"numDeleted"`
	Username   string `json:"username"`
}

// DeleteUser takes a given userID and removes the user from the database.
func DeleteUser(ctx context.Context, collection *mongo.Collection, userID string) (*mongo.DeleteResult, error) {
	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	result, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
