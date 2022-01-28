package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/team"
	"github.com/conalli/bookshelf-backend/models/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateNewUser checks whether a username alreadys exists in the db. If not, a new user
// is created based upon the request data.
func CreateNewUser(reqCtx context.Context, requestData requests.CredentialsRequest) (string, string, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	userExists := user.UserFieldAlreadyExists(ctx, &collection, "name", requestData.Name)

	if userExists {
		return "", "", errors.NewBadRequestError(fmt.Sprintf("error creating new user; user with name %v already exists", requestData.Name))
	}
	apiKey := user.GenerateAPIKey()
	for user.UserFieldAlreadyExists(ctx, &collection, "apiKey", apiKey) {
		apiKey = user.GenerateAPIKey()
	}
	hashedPassword, err := password.HashPassword(requestData.Password)
	if err != nil {
		log.Println("error hashing password")
		return "", "", errors.NewInternalServerError()
	}
	newUserData := user.NewUserData{
		Name:      requestData.Name,
		Password:  hashedPassword,
		APIKey:    apiKey,
		Bookmarks: map[string]string{},
	}
	res, err := collection.InsertOne(ctx, newUserData)
	if err != nil {
		log.Printf("error creating new user with data: \n username: %v\n password: %v", requestData.Name, requestData.Password)
		return "", "", errors.NewInternalServerError()
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("error getting objectID from newly inserted user")
		return "", "", errors.NewInternalServerError()
	}
	return oid.Hex(), apiKey, nil
}

// CreateNewTeam checks whether a team name alreadys exists in the db. If not, a new team
// is created based upon the request data.
func CreateNewTeam(reqContext context.Context, requestData requests.NewTeamRequest) (string, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqContext)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("teams")
	teamExists := user.UserFieldAlreadyExists(ctx, &collection, "name", requestData.Name)
	if teamExists {
		log.Println("team already exists")
		return "", errors.NewBadRequestError(fmt.Sprintf("error creating new user; user with name %v already exists", requestData.Name))
	}
	hashedPassword, err := password.HashPassword(requestData.Password)
	if err != nil {
		log.Printf("couldnt hash password: %+v\n", err)
		return "", errors.NewInternalServerError()
	}
	newTeamData := team.NewTeamData{
		Name:      requestData.Name,
		Password:  hashedPassword,
		ShortName: requestData.ShortName,
		Members:   map[string]string{requestData.ID: "admin"},
		Bookmarks: map[string]string{},
	}
	res, err := collection.InsertOne(ctx, newTeamData)
	if err != nil {
		log.Printf("couldnt insert team %+v\n", err)
		return "", errors.NewInternalServerError()
	}
	teamOID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("couldnt convert InsertedID to ObjectID")
		return "", errors.NewInternalServerError()
	}
	return teamOID.Hex(), nil
}
