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
)

// New checks whether a team name alreadys exists in the db. If not, a new team
// is created based upon the request data.
func New(reqContext context.Context, requestData requests.NewTeamRequest) (string, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqContext)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("teams")
	teamExists := user.DataAlreadyExists(ctx, &collection, "name", requestData.Name)
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
