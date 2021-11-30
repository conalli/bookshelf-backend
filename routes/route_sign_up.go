package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/auth/jwtauth"
	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
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
	var token string
	token, err = jwtauth.NewToken(newUserReq.Name)
	if err != nil {
		log.Printf("error returned while trying to create a new token: %v", err)
		apiErrors.APIErrorResponse(w, err)
		return
	}
	cookie := http.Cookie{
		Name:     "bookshelfjwt",
		Value:    token,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}
	log.Println("successfully returned token as cookie")
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := struct {
		Status string `json:"status"`
	}{
		Status: "success",
	}
	json.NewEncoder(w).Encode(res)
	return
}
