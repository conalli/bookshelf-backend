package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
)

// SignUp is the handler for the signup endpoint. Checks db for username and if
// unique adds new user with given credentials.
func SignUp(w http.ResponseWriter, r *http.Request) {
	log.Println("SignUp endpoint hit")
	var newUserReq models.Credentials
	json.NewDecoder(r.Body).Decode(&newUserReq)
	// add validation for request
	createUser, err := controllers.CreateNewUser(newUserReq)
	if err != nil {
		log.Printf("error returned while trying to create a new user: %v", err)
		apiErrors.APIErrorResponse(w, err)
		return
	}
	log.Printf("successfully created a new user: %v", createUser.InsertedID)
	w.WriteHeader(http.StatusCreated)
	return
}
