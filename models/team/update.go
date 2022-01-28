package team

import (
	"context"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/user"
)

// AddMember checks whether a username alreadys exists in the db. If not, a new user
// is created based upon the request data.
func AddMember(reqCtx context.Context, requestData requests.AddMemberRequest) (bool, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	userCollection := client.MongoCollection("users")
	teamCollection := client.MongoCollection("teams")
	user, err := user.GetUserByKey(ctx, &userCollection, "name", requestData.MemberName)
	if err != nil {
		return false, errors.NewBadRequestError("couldnt find user with name " + requestData.MemberName)
	}

	ok, err := AddMemberToTeam(ctx, &teamCollection, requestData.TeamID, user.ID, requestData.Role)
	if err != nil {
		log.Printf("error adding member with name: %s to team with id: %s\n error: %+v\n", requestData.MemberName, requestData.TeamID, err)
		return false, errors.NewInternalServerError()
	}
	return ok, nil
}
