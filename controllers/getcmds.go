package controllers

import (
	"context"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// GetAllCmds uses req info to get all users current cmds from the db.
func GetAllCmds(reqCtx context.Context, apiKey string) (map[string]string, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	defer cancelFunc()

	cache := db.NewRedisClient()
	cmds, err := cache.GetCachedCmds(ctx, apiKey)
	if err != nil {
		client := db.NewMongoClient(ctx)
		defer client.DB.Disconnect(ctx)

		collection := client.MongoCollection("users")
		user, err := models.GetUserByKey(ctx, &collection, "apiKey", apiKey)
		if err != nil {
			return nil, apiErrors.ParseGetUserError(apiKey, err)
		}
		cache.SetCacheCmds(ctx, apiKey, user.Bookmarks)
		return user.Bookmarks, nil
	}
	return cmds, nil
}