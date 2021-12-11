package controllers

import (
	"context"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
	"go.mongodb.org/mongo-driver/mongo"
)

// DelCmd attempts to either rempve a cmd from the user, returning the number
// of updated cmds.
func DelCmd(reqCtx context.Context, requestData models.DelCmdReq) (int, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	user, err := models.GetUserByKey(ctx, &collection, "name", requestData.Name)
	if err != nil {
		return 0, apiErrors.ParseGetUserError(requestData.Name, err)
	}
	correctPassword := password.CheckHashedPassword(user.Password, requestData.Password)
	if !correctPassword {
		return 0, apiErrors.NewWrongCredentialsError("error: password incorrect")
	}
	var result *mongo.UpdateResult
	result, err = models.RemoveCmdFromUser(ctx, &collection, user.Name, requestData.Cmd)
	if err != nil {
		return 0, apiErrors.NewInternalServerError()
	}
	return int(result.ModifiedCount), nil
}
