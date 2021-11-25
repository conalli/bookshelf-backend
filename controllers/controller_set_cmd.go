package controllers

import (
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
	"github.com/conalli/bookshelf-backend/utils/password"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddCmd attempts to either add or update a cmd for the user, returning the number
// of updated cmds.
func AddCmd(requestData models.SetCmdReq) (int, apiErrors.ApiErr) {
	ctx, cancelFunc := db.MongoContext()
	client := db.MongoClient(ctx)
	defer cancelFunc()
	defer client.Disconnect(ctx)

	collection := db.MongoCollection(client, "users")
	user, err := models.GetUserByKey(ctx, &collection, "name", requestData.Name)
	if err != nil {
		return 0, apiErrors.ParseGetUserError(requestData.Name, err)
	}
	correctPassword := password.CheckHashedPassword(user.Password, requestData.Password)
	if !correctPassword {
		return 0, apiErrors.NewWrongCredentialsError("error: password incorrect")
	}
	var result *mongo.UpdateResult
	result, err = models.AddCmdToUser(ctx, &collection, user.Name, requestData.Cmd, requestData.URL)
	if err != nil {
		return 0, apiErrors.NewInternalServerError()
	}
	var numUpdated int
	if int(result.UpsertedCount) >= int(result.ModifiedCount) {
		numUpdated = int(result.UpsertedCount)
	} else {
		numUpdated = int(result.ModifiedCount)
	}
	return numUpdated, nil
}
