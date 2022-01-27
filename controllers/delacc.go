package controllers

import (
	"context"
	"log"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// DelAcc attempts to delete a user from the db, returning the number of deleted users.
func DelAcc(reqCtx context.Context, requestData models.DelAccRequest, apiKey string) (int, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	user, err := models.GetUserByID(ctx, &collection, requestData.ID)
	if err != nil {
		log.Printf("error deleting user: couldn't find user -> %v", err)
		return 0, apiErrors.NewBadRequestError("could not find user to delete")
	}
	ok := password.CheckHashedPassword(user.Password, requestData.Password)
	if !ok {
		log.Printf("error deleting user: password incorrect -> %v", err)
		return 0, apiErrors.NewWrongCredentialsError("password incorrect")
	}
	result, err := models.DeleteUser(ctx, &collection, requestData.ID)
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
