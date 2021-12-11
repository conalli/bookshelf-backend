package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateNewUser checks whether a username alreadys exists in the db. If not, a new user
// is created based upon the request data.
func CreateNewUser(reqCtx context.Context, requestData models.Credentials) (*mongo.InsertOneResult, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	userExists := models.UserFieldAlreadyExists(ctx, &collection, "name", requestData.Name)

	if !userExists {
		apiKey := models.GenerateAPIKey()
		for models.UserFieldAlreadyExists(ctx, &collection, "apiKey", apiKey) {
			apiKey = models.GenerateAPIKey()
		}
		hashedPassword, hashErr := password.HashPassword(requestData.Password)
		if hashErr != nil {
			log.Println("error hashing password")
			return nil, apiErrors.NewInternalServerError()
		}
		newUserData := models.UserData{
			Name:      requestData.Name,
			Password:  hashedPassword,
			APIKey:    apiKey,
			Bookmarks: map[string]string{},
		}
		res, err := collection.InsertOne(ctx, newUserData)
		if err != nil {
			log.Printf("error creating new user with data: \n username: %v\n password: %v", requestData.Name, requestData.Password)
			return nil, apiErrors.NewInternalServerError()
		}
		return res, nil
	}
	return nil, apiErrors.NewBadRequestError(fmt.Sprintf("error creating new user; user with name %v already exists", requestData.Name))
}
