package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserData represents the document fields created when a user signs up.
type NewUserData struct {
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	APIKey    string            `json:"apiKey" bson:"apiKey"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}

// UserData represents the db fields associated with each user.
type UserData struct {
	ID        string            `json:"id" bson:"_id"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	APIKey    string            `json:"apiKey" bson:"apiKey"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}

// GetUserByID finds and returns user data based on a the users _id.
func GetUserByID(ctx context.Context, collection *mongo.Collection, userID string) (UserData, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return UserData{}, err
	}
	var result UserData
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, mongo.ErrNoDocuments
		}
		return result, err
	}
	return result, nil
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
