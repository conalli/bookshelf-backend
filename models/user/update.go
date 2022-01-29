package user

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateEmbedOptions struct {
	FilterKey, FilterValue, Embedded, Key, Value string
}

func UpdateEmbedByField(ctx context.Context, collection *mongo.Collection, data UpdateEmbedOptions) (User, error) {
	options := options.FindOneAndUpdate().SetUpsert(true)
	var filter primitive.M
	if data.FilterKey == "_id" {
		userID, err := primitive.ObjectIDFromHex(data.FilterValue)
		if err != nil {
			return User{}, err
		}
		filter = bson.M{data.FilterKey: userID}
	} else {
		filter = bson.M{data.FilterKey: data.FilterValue}
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("%s.%s", data.Embedded, data.Key), Value: data.Value}}}}
	var result User
	err := collection.FindOneAndUpdate(ctx, filter, update, options).Decode(&result)
	if err != nil {
		return User{}, err
	}
	return result, nil
}

// AddCmd attempts to either add or update a cmd for the user, returning the number
// of updated cmds.
func AddCmd(reqCtx context.Context, requestData requests.AddCmdRequest, apiKey string) (int, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")

	result, err := AddCmdToUser(ctx, &collection, requestData.ID, requestData.Cmd, requestData.URL)
	if err != nil {
		return 0, errors.NewInternalServerError()
	}
	var numUpdated int
	if int(result.UpsertedCount) >= int(result.ModifiedCount) {
		numUpdated = int(result.UpsertedCount)
	} else {
		numUpdated = int(result.ModifiedCount)
	}
	if numUpdated >= 1 {
		cache := db.NewRedisClient()
		cmds, err := cache.GetCachedCmds(ctx, apiKey)
		if err != nil {
			log.Println("could not get cached cmds after adding new cmd")
		} else {
			cmds[requestData.Cmd] = requestData.URL
			cache.SetCacheCmds(ctx, apiKey, cmds)
			log.Printf("successfully updated cache with new cmd %s:%s\n", requestData.Cmd, requestData.URL)
		}
	}
	return numUpdated, nil
}

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

// DelCmd attempts to either rempve a cmd from the user, returning the number
// of updated cmds.
func DelCmd(reqCtx context.Context, requestData requests.DelCmdRequest, apiKey string) (int, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	result, err := RemoveCmdFromUser(ctx, &collection, requestData.ID, requestData.Cmd)
	if err != nil {
		return 0, errors.NewInternalServerError()
	}
	if result.ModifiedCount >= 1 {
		cache := db.NewRedisClient()
		cmds, err := cache.GetCachedCmds(ctx, apiKey)
		if err != nil {
			log.Println("could not get cached cmds after removing cmd")
		}
		delete(cmds, requestData.Cmd)
		cache.SetCacheCmds(ctx, apiKey, cmds)
		log.Printf("successfully removed cmd: %s from cache \n", requestData.Cmd)
	}
	return int(result.ModifiedCount), nil
}

// RemoveCmdFromUser takes a given username along with the cmd and removes the cmd from their bookmarks.
func RemoveCmdFromUser(ctx context.Context, collection *mongo.Collection, userID, cmd string) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	update := bson.D{primitive.E{Key: "$unset", Value: bson.D{primitive.E{Key: fmt.Sprintf("bookmarks.%s", cmd), Value: ""}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
