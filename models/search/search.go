package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/user"
)

// GetURL takes in an apiKey and cmd and returns either a correctly formatted url from the db,
// or a google search url for the cmd based on whether the cmd could be found or not.
func GetURL(reqCtx context.Context, apiKey, cmd string) (string, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	defer cancelFunc()

	cache := db.NewRedisClient()
	url, err := cache.GetSearchData(ctx, apiKey, cmd)
	if err != nil {
		client := db.NewMongoClient(ctx)
		defer client.DB.Disconnect(ctx)
		collection := client.MongoCollection("users")

		currUser, err := user.GetByKey(ctx, &collection, "apiKey", apiKey)
		defaultSearch := fmt.Sprintf("http://www.google.com/search?q=%s", cmd)
		if err != nil {
			return defaultSearch, errors.ParseGetUserError(apiKey, err)
		}

		cache.SetCacheCmds(ctx, apiKey, currUser.Bookmarks)

		url, found := currUser.Bookmarks[cmd]
		if !found {
			return defaultSearch, errors.NewBadRequestError("error: command: " + cmd + " not registered")
		}
		return FormatURL(url), nil
	}
	return FormatURL(url), nil
}

// FormatURL takes a url from the database and returns it in a format required
// for successful redirection.
func FormatURL(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}
	return "http://" + url
}
