package search

import (
	"strings"
)

// // GetURL takes in an APIKey and cmd and returns either a correctly formatted url from the db,
// // or a google search url for the cmd based on whether the cmd could be found or not.
// func GetURL(reqCtx context.Context, APIKey, cmd string) (string, errors.ApiErr) {
// 	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
// 	defer cancelFunc()

// 	cache := redis.NewClient()
// 	url, err := cache.GetSearchData(ctx, APIKey, cmd)
// 	if err != nil {

// 	}
// 	return formatURL(url), nil
// }

func formatURL(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}
	return "http://" + url
}
