package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SetCmdReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Cmd      string `json:"cmd" bson:"cmd"`
	URL      string `json:"url" bson:"url"`
}

type SetCmdRes struct {
	CmdsSet int `json:"cmdsSet"`
}

func AddCmdToUser(ctx context.Context, collection *mongo.Collection, username, key, value string) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"name", username}}
	update := bson.D{{"$set", bson.D{{fmt.Sprintf("bookmarks.%s", key), value}}}}
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
