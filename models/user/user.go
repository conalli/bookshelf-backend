package user

import (
	"context"
	"net/http"

	"github.com/conalli/bookshelf-backend/auth/password"
	"github.com/conalli/bookshelf-backend/db"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
)

// NewUserData represents the document fields created when a user signs up.
type NewUserData struct {
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	APIKey    string            `json:"apiKey" bson:"apiKey"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
	Teams     map[string]string `json:"teams" bson:"teams"`
}

// User represents the db fields associated with each user.
type User struct {
	ID        string            `json:"id" bson:"_id"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	APIKey    string            `json:"apiKey" bson:"apiKey"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
	Teams     map[string]string `json:"teams" bson:"teams"`
}

// CheckCredentials takes in request data, checks the db and returns the username and apikey is successful.
func CheckCredentials(reqCtx context.Context, requestData requests.CredentialsRequest) (User, errors.ApiErr) {
	ctx, cancelFunc := db.ReqContextWithTimeout(reqCtx)
	client := db.NewMongoClient(ctx)
	defer cancelFunc()
	defer client.DB.Disconnect(ctx)

	collection := client.MongoCollection("users")
	currUser, err := GetByKey(ctx, &collection, "name", requestData.Name)
	if err != nil || !password.CheckHashedPassword(currUser.Password, requestData.Password) {
		return User{}, errors.NewApiError(http.StatusUnauthorized, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	return currUser, nil
}
