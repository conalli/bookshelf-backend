package user

import (
	"context"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetByID finds and returns user data based on a the users _id.
func GetByID(ctx context.Context, collection *mongo.Collection, userID string) (UserData, error) {
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

// GetByKey finds and returns user data based on a key-value pair.
func GetByKey(ctx context.Context, collection *mongo.Collection, reqKey, reqValue string) (UserData, error) {
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

// GetAllCmds uses req info to get all users current cmds from the db.
func GetAllCmds(reqCtx context.Context, apiKey string) (map[string]string, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	defer cancelFunc()

	cache := db.NewRedisClient()
	cmds, err := cache.GetCachedCmds(ctx, apiKey)
	if err != nil {
		client := db.NewMongoClient(ctx)
		defer client.DB.Disconnect(ctx)

		collection := client.MongoCollection("users")
		user, err := GetByKey(ctx, &collection, "apiKey", apiKey)
		if err != nil {
			return nil, errors.ParseGetUserError(apiKey, err)
		}
		cache.SetCacheCmds(ctx, apiKey, user.Bookmarks)
		return user.Bookmarks, nil
	}
	return cmds, nil
}
