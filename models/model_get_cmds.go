package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetCmdsReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func GetUserByName(ctx context.Context, collection *mongo.Collection, reqName string) (UserData, error) {
	var result UserData
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: "name", Value: reqName}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, mongo.ErrNoDocuments
		} else {
			return result, err
		}
	}
	return result, nil
}
