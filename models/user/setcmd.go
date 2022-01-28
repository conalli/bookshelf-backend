package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddCmdToUser takes a given username along with the cmd and URL to set and adds the data to their bookmarks.
func AddCmdToUser(ctx context.Context, collection *mongo.Collection, userID, key, value string) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("bookmarks.%s", key), Value: value}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
