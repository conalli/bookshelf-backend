package controllers

import (
	"context"
	"log"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// DelCmd attempts to either rempve a cmd from the user, returning the number
// of updated cmds.
func DelCmd(reqCtx context.Context, requestData models.DelCmdReq, apiKey string) (int, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	result, err := models.RemoveCmdFromUser(ctx, &collection, requestData.ID, requestData.Cmd)
	if err != nil {
		return 0, apiErrors.NewInternalServerError()
	}
	if result.ModifiedCount >= 1 {
		cache := db.NewRedisClient()
		cmds, err := cache.GetCachedCmds(ctx, apiKey)
		if err != nil {
			log.Println("could not get cached cmds after removing cmd")
		}
		delete(cmds, requestData.Cmd)
		cache.SetCacheCmds(ctx, apiKey, cmds)
		log.Printf("successfully removed cmd: %s from cache \n", requestData.Cmd)
	}
	return int(result.ModifiedCount), nil
}