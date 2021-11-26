package controllers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
	"github.com/conalli/bookshelf-backend/utils/auth/password"
)

func CheckCredentials(requestData models.Credentials) (string, apiErrors.ApiErr) {
	ctx, cancelFunc := db.MongoContext()
	client := db.MongoClient(ctx)
	defer cancelFunc()
	defer client.Disconnect(ctx)

	collection := db.MongoCollection(client, "users")
	user, err := models.GetUserByKey(ctx, &collection, "name", requestData.Name)
	if err != nil || !password.CheckHashedPassword(user.Password, requestData.Password) {
		return "", apiErrors.NewApiError(http.StatusUnauthorized, apiErrors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	return user.Name, nil
}
