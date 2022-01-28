package user

import (
	"context"
	"log"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Delete attempts to delete a user from the db, returning the number of deleted users.
func Delete(reqCtx context.Context, requestData requests.DelUserRequest, apiKey string) (int, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	userData, err := GetByID(ctx, &collection, requestData.ID)
	if err != nil {
		log.Printf("error deleting user: couldn't find user -> %v", err)
		return 0, errors.NewBadRequestError("could not find user to delete")
	}
	ok := password.CheckHashedPassword(userData.Password, requestData.Password)
	if !ok {
		log.Printf("error deleting user: password incorrect -> %v", err)
		return 0, errors.NewWrongCredentialsError("password incorrect")
	}
	result, err := DeleteUserFromDB(ctx, &collection, requestData.ID)
	if err != nil {
		log.Printf("error deleting user: error -> %v", err)
		return 0, errors.NewInternalServerError()
	}
	if result.DeletedCount == 0 {
		log.Printf("could not remove user... maybe user:%s doesn't exists?", requestData.Name)
		return 0, errors.NewBadRequestError("error: could not remove cmd")
	}
	cache := db.NewRedisClient()
	cache.DelCachedCmds(ctx, apiKey)
	return int(result.DeletedCount), nil
}

// DeleteUserFromDB takes a given userID and removes the user from the database.
func DeleteUserFromDB(ctx context.Context, collection *mongo.Collection, userID string) (*mongo.DeleteResult, error) {
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
