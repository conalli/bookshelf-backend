package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
)

type tempError struct {
	Error string `json:"error"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	log.Println("SignUp endpoint hit")
	var newUserReq models.SignUpReq
	json.NewDecoder(r.Body).Decode(&newUserReq)
	// add validation for request
	createUser, err := controllers.CreateNewUser(newUserReq)
	if err != nil {
		log.Printf("error returned while trying to create a new user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		var testError tempError
		testError.Error = err.Error()
		json.NewEncoder(w).Encode(testError)
	} else {
		log.Printf("successfully created a new user: %v", createUser.InsertedID)
		w.WriteHeader(http.StatusCreated)
	}
}
