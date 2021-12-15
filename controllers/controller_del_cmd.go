package controllers

import (
	"context"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// DelCmd attempts to either rempve a cmd from the user, returning the number
// of updated cmds.
func DelCmd(reqCtx context.Context, requestData models.DelCmdReq) (int, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	result, err := models.RemoveCmdFromUser(ctx, &collection, requestData.ID, requestData.Cmd)
	if err != nil {
		return 0, apiErrors.NewInternalServerError()
	}
	return int(result.ModifiedCount), nil
}
