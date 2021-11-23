package controllers

import (
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
	"github.com/conalli/bookshelf-backend/utils/password"
)

func GetAllCmds(requestData models.GetCmdsReq) (map[string]string, error) {
	ctx, cancelFunc := db.MongoContext()
	client := db.MongoClient(ctx)
	defer cancelFunc()
	defer client.Disconnect(ctx)

	collection := db.MongoCollection(client, "users")
	user, err := models.GetUserByKey(ctx, &collection, "name", requestData.Name)
	if err != nil {
		return nil, apiErrors.ParseGetUserError(requestData.Name, err)
	}
	correctPassword := password.CheckHashedPassword(user.Password, requestData.Password)
	if !correctPassword {
		return nil, apiErrors.NewWrongCredentialsError("error: password incorrect")
	}
	return user.Bookmarks, nil
}
