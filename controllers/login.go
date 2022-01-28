package controllers

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

// LogIn is the handler for the login endpoint. Checks credentials and if
// correct returns JWT cookie for use with getcmds and setcmd.
func LogIn(w http.ResponseWriter, r *http.Request) {
	log.Println("LogIn endpoint hit")
	var logInReq requests.CredentialsRequest
	json.NewDecoder(r.Body).Decode(&logInReq)

	currUser, err := user.CheckCredentials(r.Context(), logInReq)
	if err != nil {
		log.Printf("error returned while trying to get check credentials: %v", err)
		errors.APIErrorResponse(w, err)
		return
	}
	token, err := jwtauth.NewToken(currUser.APIKey)
	if err != nil {
		log.Printf("error returned while trying to create a new token: %v", err)
		errors.APIErrorResponse(w, err)
		return
	}
	// Use Secure during production.
	cookie := http.Cookie{
		Name:     "bookshelfjwt",
		Value:    token,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	log.Println("successfully returned token as cookie")
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := responses.LogInResponse{
		ID:     currUser.ID,
		APIKey: currUser.APIKey,
	}
	json.NewEncoder(w).Encode(res)
	return
}
