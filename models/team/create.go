package team

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// New checks whether a team name alreadys exists in the db. If not, a new team
// is created based upon the request data.
func New(reqContext context.Context, requestData requests.NewTeamRequest) (string, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqContext)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	teamCollection := client.MongoCollection("teams")
	teamExists := user.DataAlreadyExists(ctx, &teamCollection, "name", requestData.Name)
	if teamExists {
		log.Println("team already exists")
		return "", errors.NewBadRequestError(fmt.Sprintf("error creating new user; user with name %v already exists", requestData.Name))
	}
	hashedPassword, err := password.HashPassword(requestData.Password)
	if err != nil {
		log.Printf("couldnt hash password: %+v\n", err)
		return "", errors.NewInternalServerError()
	}
	newTeamData := NewTeamData{
		Name:      requestData.Name,
		Password:  hashedPassword,
		ShortName: requestData.ShortName,
		Members:   map[string]string{requestData.ID: "admin"},
		Bookmarks: map[string]string{},
	}
	res, err := client.SessionWithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		result, err := teamCollection.InsertOne(sessCtx, newTeamData)
		if err != nil {
			log.Printf("couldnt insert team %+v\n", err)
			return "", errors.NewInternalServerError()
		}
		teamOID, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Println("couldnt convert InsertedID to ObjectID")
			return "", errors.NewInternalServerError()
		}
		userCollection := client.MongoCollection("users")
		opts := user.UpdateEmbedOptions{
			FilterKey:   "_id",
			FilterValue: requestData.ID,
			Embedded:    "teams",
			Key:         teamOID.Hex(),
			Value:       "admin",
		}
		user, err := user.UpdateEmbedByField(sessCtx, &userCollection, opts)
		if err != nil {
			log.Printf("couldnt add team: %s to user: %s", newTeamData.Name, user.Name)
		}
		return teamOID.Hex(), nil
	})
	if err != nil {
		log.Println("error during session transaction")
		return "", errors.NewInternalServerError()
	}
	v, ok := res.(string)
	if !ok {
		return "", errors.NewInternalServerError()
	}
	return v, nil
}
