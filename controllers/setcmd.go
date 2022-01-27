package controllers

import (
	"context"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// AddCmd attempts to either add or update a cmd for the user, returning the number
// of updated cmds.
func AddCmd(reqCtx context.Context, requestData models.SetCmdReq, apiKey string) (int, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")

	result, err := models.AddCmdToUser(ctx, &collection, requestData.ID, requestData.Cmd, requestData.URL)
	if err != nil {
		return 0, apiErrors.NewInternalServerError()
	}
	var numUpdated int
	if int(result.UpsertedCount) >= int(result.ModifiedCount) {
		numUpdated = int(result.UpsertedCount)
	} else {
		numUpdated = int(result.ModifiedCount)
	}
	if numUpdated >= 1 {
		cache := db.NewRedisClient()
		cmds, err := cache.GetCachedCmds(ctx, apiKey)
		if err != nil {
			log.Println("could not get cached cmds after adding new cmd")
		} else {
			cmds[requestData.Cmd] = requestData.URL
			cache.SetCacheCmds(ctx, apiKey, cmds)
			log.Printf("successfully updated cache with new cmd %s:%s\n", requestData.Cmd, requestData.URL)
		}
	}
	return numUpdated, nil
}