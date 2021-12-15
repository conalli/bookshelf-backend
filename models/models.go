package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Credentials represents the fields needed in the request in order to attempt to sign up or log in.
type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// UserData represents the db fields associated with each user.
type UserData struct {
	ID        string            `json:"id" bson:"_id"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	APIKey    string            `json:"apiKey" bson:"apiKey"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}

// SuccessRes represents a general success message to return to the user.
type SuccessRes struct {
	Status string `json:"status"`
}

// GetUserByKey finds and returns user data based on a key-value pair.
func GetUserByKey(ctx context.Context, collection *mongo.Collection, reqKey, reqValue string) (UserData, error) {
	var result UserData
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: reqKey, Value: reqValue}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, mongo.ErrNoDocuments
		}
		return result, err
	}
	return result, nil
}
