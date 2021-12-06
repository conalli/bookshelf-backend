package controllers

import (
	"fmt"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// GetURL takes in an apiKey and cmd and returns either a correctly formatted url from the db,
// or a google search url for the cmd based on whether the cmd could be found or not.
func GetURL(apiKey, cmd string) (string, apiErrors.ApiErr) {
	ctx, cancelFunc := db.MongoContext()
	defer cancelFunc()

	cache := db.NewRedisClient()
	url, err := cache.GetSearchData(ctx, apiKey, cmd)
	if err != nil {
		client := db.MongoClient(ctx)
		defer client.Disconnect(ctx)
		collection := db.MongoCollection(client, "users")

		var user models.UserData
		user, err = models.GetUserByKey(ctx, &collection, "apiKey", apiKey)
		var defaultSearch = fmt.Sprintf("http://www.google.com/search?q=%s", cmd)
		if err != nil {
			return defaultSearch, apiErrors.ParseGetUserError(apiKey, err)
		}

		cache.SetSearchData(ctx, apiKey, user.Bookmarks)

		url, found := user.Bookmarks[cmd]
		if !found {
			return defaultSearch, apiErrors.NewBadRequestError("error: command: " + cmd + " not registered")
		}
		return models.FormatURL(url), nil
	}
	return models.FormatURL(url), nil
}
