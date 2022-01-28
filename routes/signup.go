package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/auth/jwtauth"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/responses"
	"github.com/conalli/bookshelf-backend/models/user"
)

// SignUp is the handler for the signup endpoint. Checks db for username and if
// unique adds new user with given credentials.
func SignUp(w http.ResponseWriter, r *http.Request) {
	log.Println("SignUp endpoint hit")
	var newUserReq requests.CredentialsRequest
	json.NewDecoder(r.Body).Decode(&newUserReq)
	userID, apiKey, err := user.New(r.Context(), newUserReq)
	if err != nil {
		log.Printf("error returned while trying to create a new user: %v", err)
		errors.APIErrorResponse(w, err)
		return
	}
	log.Printf("successfully created a new user: %s", userID)
	var token string
	token, err = jwtauth.NewToken(apiKey)
	if err != nil {
		log.Printf("error returned while trying to create a new token: %v", err)
		errors.APIErrorResponse(w, err)
		return
	}
	cookie := http.Cookie{
		Name:     "bookshelfjwt",
		Value:    token,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	log.Println("successfully returned token as cookie")
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := responses.SignUpResponse{
		ID:     userID,
		APIKey: apiKey,
	}
	json.NewEncoder(w).Encode(res)
	return
}
