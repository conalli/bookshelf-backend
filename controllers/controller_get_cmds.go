package controllers

import (
	"fmt"

	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models"
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
		return nil, err
	}
	correctPassword := password.CheckHashedPassword(user.Password, requestData.Password)
	if !correctPassword {
		return nil, fmt.Errorf("could not return commands, incorrect password")
	}
	return user.Bookmarks, nil
}
