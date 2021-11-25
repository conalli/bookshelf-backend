package controllers

import (
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
	"github.com/conalli/bookshelf-backend/utils/password"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateNewUser(requestData models.SignUpReq) (*mongo.InsertOneResult, apiErrors.ApiErr) {
	ctx, cancelFunc := db.MongoContext()
	client := db.MongoClient(ctx)
	defer cancelFunc()
	defer client.Disconnect(ctx)

	collection := db.MongoCollection(client, "users")
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
