package controllers

import (
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateNewUser(requestData models.SignUpReq) (*mongo.InsertOneResult, error) {
	ctx, cancelFunc := db.MongoContext()
	client := db.MongoClient(ctx)
	defer cancelFunc()
	defer client.Disconnect(ctx)

	collection := db.MongoCollection(client, "users")
	userExists := models.UserAlreadyExists(ctx, &collection, *requestData.Name)

	if !userExists {
		apiKey := models.GenerateAPIKey()
		newUserData := models.SignUpData{
			Name:     *requestData.Name,
			Password: *requestData.Password,
			ApiKey:   apiKey,
		}
		res, err := collection.InsertOne(ctx, newUserData)
		if err != nil {
			log.Printf("error creating new user with data: \n username: %v\n password: %v", requestData.Name, requestData.Password)
		}
		return res, err
	}
	return nil, fmt.Errorf("error creating new user; user with name %v already exists", requestData.Name)
}
