package team

import (
	"context"
	"fmt"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddMember checks whether a username alreadys exists in the db. If not, a new user
// is created based upon the request data.
func AddMember(reqCtx context.Context, requestData requests.AddMemberRequest) (bool, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)
	res, err := client.SessionWithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		userCollection := client.MongoCollection("users")
		update := user.UpdateEmbedOptions{
			FilterKey:   "name",
			FilterValue: requestData.MemberName,
			Embedded:    "teams",
			Key:         requestData.TeamID,
			Value:       requestData.Role,
		}
		user, err := user.UpdateEmbedByField(sessCtx, &userCollection, update)
		if err != nil {
			log.Println("couldnt update user with name: " + requestData.MemberName)
			return false, errors.NewBadRequestError("couldnt update user with name: " + requestData.MemberName)
		}
		teamCollection := client.MongoCollection("teams")
		ok, err := AddMemberToTeam(ctx, &teamCollection, requestData.TeamID, user.ID, requestData.Role)
		if err != nil {
			log.Printf("error adding member with name: %s to team with id: %s\n error: %+v\n", requestData.MemberName, requestData.TeamID, err)
			return false, errors.NewInternalServerError()
		}
		return ok, nil
	})
	if err != nil {
		log.Println("error could not start db transaction")
		return false, errors.NewInternalServerError()
	}
	if v, ok := res.(bool); ok {
		if !v {
			return false, nil
		}
	} else {
		return false, nil
	}
	return true, nil
}

// AddMemberToTeam attempts to add a new member to a team.
func AddMemberToTeam(ctx context.Context, collection *mongo.Collection, teamID, memberID, role string) (bool, error) {
	opts := options.Update().SetUpsert(true)
	filter, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		log.Printf("error getting objectid from hex: %+v\n", err)
		return false, err
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: fmt.Sprintf("members.%s", memberID), Value: role}}}}
	result, err := collection.UpdateByID(ctx, filter, update, opts)
	if err != nil {
		log.Printf("error attempting to add user: %s to team: %s: %+v\n", memberID, teamID, err)
		return false, err
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		log.Printf("error attempting to add user: %s to team: %s, team was not modified\n", memberID, teamID)
		return false, nil
	}
	return true, nil
}
