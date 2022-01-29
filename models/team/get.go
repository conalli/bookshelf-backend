package team

import (
	"context"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetTeams uses user id to get all users teams from the db.
func GetTeams(ctx context.Context, apiKey string) ([]Team, errors.ApiErr) {
	reqCtx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	client := db.NewMongoClient(reqCtx)
	defer client.DB.Disconnect(reqCtx)
	res, err := client.SessionWithTransaction(reqCtx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		userCollection := client.MongoCollection("users")
		user, err := user.GetByKey(sessCtx, &userCollection, "apiKey", apiKey)
		if err != nil {
			log.Printf("error getting user by apiKey: %s -> %+v\n", apiKey, err)
			return nil, err
		}
		teamIDs, err := convertIDs(user.Teams)
		if err != nil {
			log.Printf("error converting teams to ids: error -> %+v\n", err)
			return nil, err
		}
		teamCollection := client.MongoCollection("teams")
		filter := bson.M{"_id": bson.M{"$in": teamIDs}}
		opts := options.Find()
		teamCursor, err := teamCollection.Find(sessCtx, filter, opts)
		defer teamCursor.Close(sessCtx)
		if err != nil {
			log.Printf("error converting teams to ids -> %+v\n", err)
			return nil, err
		}
		var teams []Team
		for teamCursor.Next(sessCtx) {
			var currTeam Team
			if err := teamCursor.Decode(&currTeam); err != nil {
				log.Printf("error could not get team from found teams -> %+v\n", err)
				return nil, err
			}
			teams = append(teams, currTeam)
		}
		return teams, nil
	})
	if err != nil {
		log.Printf("error could not get data from transaction -> %+v\n", err)
		return nil, errors.NewInternalServerError()
	}
	teams, ok := res.([]Team)
	if !ok {
		log.Println("error could not assert type []Team")
		return nil, errors.NewInternalServerError()
	}
	return teams, nil
}

func convertIDs(teams map[string]string) ([]primitive.ObjectID, error) {
	output := make([]primitive.ObjectID, len(teams))
	for id := range teams {
		res, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		output = append(output, res)
	}
	return output, nil
}
