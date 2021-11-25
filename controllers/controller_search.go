package controllers

import (
	"fmt"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
)

func GetURL(apiKey, cmd string) (string, apiErrors.ApiErr) {
	ctx, cancelFunc := db.MongoContext()
	client := db.MongoClient(ctx)
	defer cancelFunc()
	defer client.Disconnect(ctx)

	collection := db.MongoCollection(client, "users")
	user, err := models.GetUserByKey(ctx, &collection, "apiKey", apiKey)
	var defaultSearch = fmt.Sprintf("http://www.google.com/search?q=%s", cmd)
	if err != nil {
		return defaultSearch, apiErrors.ParseGetUserError(apiKey, err)
	}
	url, found := user.Bookmarks[cmd]
	if !found {
		return defaultSearch, apiErrors.NewBadRequestError("error: command: " + cmd + " not registered")
	}
	url = models.FormatURL(url)
	return url, nil
}
