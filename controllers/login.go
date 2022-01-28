package controllers

import (
	"context"
	"net/http"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/user"
)

// CheckCredentials takes in request data, checks the db and returns the username and apikey is successful.
func CheckCredentials(reqCtx context.Context, requestData requests.CredentialsRequest) (user.UserData, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	currUser, err := user.GetUserByKey(ctx, &collection, "name", requestData.Name)
	if err != nil || !password.CheckHashedPassword(currUser.Password, requestData.Password) {
		return user.UserData{}, errors.NewApiError(http.StatusUnauthorized, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	return currUser, nil
}
