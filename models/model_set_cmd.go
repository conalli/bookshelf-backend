package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetCmdReq represents the expected fields needed for the setcmd request to be completed.
type SetCmdReq struct {
	Name     string `json:"name" bson:"name"`
	Password string `json:"password"`
	Cmd      string `json:"cmd" bson:"cmd"`
	URL      string `json:"url" bson:"url"`
}

// SetCmdRes represents the number of commands that have been updated by the setcmd request.
type SetCmdRes struct {
	CmdsSet int `json:"cmdsSet"`
}

// AddCmdToUser takes a given username along with the cmd and URL to set and adds the data to their bookmarks.
func AddCmdToUser(ctx context.Context, collection *mongo.Collection, username, key, value string) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "name", Value: username}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("bookmarks.%s", key), Value: value}}}}
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
