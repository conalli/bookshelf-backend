package controllers

import (
	"context"
	"net/http"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// CheckCredentials takes in request data, checks the db and returns the username and apikey is successful.
func CheckCredentials(reqCtx context.Context, requestData models.Credentials) (string, string, apiErrors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	user, err := models.GetUserByKey(ctx, &collection, "name", requestData.Name)
	if err != nil || !password.CheckHashedPassword(user.Password, requestData.Password) {
		return "", "", apiErrors.NewApiError(http.StatusUnauthorized, apiErrors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	return user.Name, user.APIKey, nil
}
