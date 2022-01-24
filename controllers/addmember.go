package controllers

import (
	"context"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// AddMember checks whether a username alreadys exists in the db. If not, a new user
// is created based upon the request data.
func AddMember(reqCtx context.Context, requestData models.AddMemberReq) (bool, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	userCollection := client.MongoCollection("users")
	teamCollection := client.MongoCollection("teams")
	user, err := models.GetUserByKey(ctx, &userCollection, "name", requestData.MemberName)
	if err != nil {
		return false, apiErrors.NewBadRequestError("couldnt find user with name " + requestData.MemberName)
	}

	ok, err := models.AddMemberToTeam(ctx, &teamCollection, requestData.TeamID, user.ID, requestData.Role)
	if err != nil {
		log.Printf("error adding member with name: %s to team with id: %s\n error: %+v\n", requestData.MemberName, requestData.TeamID, err)
		return false, apiErrors.NewInternalServerError()
	}
	return ok, nil
}
