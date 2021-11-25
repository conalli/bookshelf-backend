package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserData represents the db fields associated with each user.
type UserData struct {
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	APIKey    string            `json:"apiKey" bson:"apiKey"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}

// GetUserByKey finds and returns user data based on a key-value pair.
func GetUserByKey(ctx context.Context, collection *mongo.Collection, reqKey, reqValue string) (UserData, error) {
	var result UserData
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: reqKey, Value: reqValue}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, mongo.ErrNoDocuments
		} else {
			return result, err
		}
	}
	return result, nil
}
