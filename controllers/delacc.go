package controllers

import (
	"context"
	"log"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/user"
)

// DelAcc attempts to delete a user from the db, returning the number of deleted users.
func DelAcc(reqCtx context.Context, requestData requests.DelAccRequest, apiKey string) (int, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	userData, err := user.GetUserByID(ctx, &collection, requestData.ID)
	if err != nil {
		log.Printf("error deleting user: couldn't find user -> %v", err)
		return 0, apiErrors.NewBadRequestError("could not find user to delete")
	}
	ok := password.CheckHashedPassword(userData.Password, requestData.Password)
	if !ok {
		log.Printf("error deleting user: password incorrect -> %v", err)
		return 0, apiErrors.NewWrongCredentialsError("password incorrect")
	}
	result, err := user.DeleteUser(ctx, &collection, requestData.ID)
	if err != nil {
		log.Printf("error deleting user: error -> %v", err)
		return 0, apiErrors.NewInternalServerError()
	}
	if result.DeletedCount == 0 {
		log.Printf("could not remove user... maybe user:%s doesn't exists?", requestData.Name)
		return 0, apiErrors.NewBadRequestError("error: could not remove cmd")
	}
	cache := db.NewRedisClient()
	cache.DelCachedCmds(ctx, apiKey)
	return int(result.DeletedCount), nil
}
