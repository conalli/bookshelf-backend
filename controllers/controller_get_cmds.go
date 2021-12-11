package controllers

import (
	"context"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// GetAllCmds uses req info to get all users current cmds from the db.
func GetAllCmds(reqCtx context.Context, userName string) (map[string]string, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	user, err := models.GetUserByKey(ctx, &collection, "name", userName)
	if err != nil {
		return nil, apiErrors.ParseGetUserError(userName, err)
	}
	return user.Bookmarks, nil
}
