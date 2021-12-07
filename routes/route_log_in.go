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

// LogIn is the handler for the login endpoint. Checks credentials and if
// correct returns JWT cookie for use with getcmds and setcmd.
func LogIn(w http.ResponseWriter, r *http.Request) {
	log.Println("LogIn endpoint hit")
	var logInReq models.Credentials
	json.NewDecoder(r.Body).Decode(&logInReq)

	name, apiKey, err := controllers.CheckCredentials(r.Context(), logInReq)
	if err != nil {
		log.Printf("error returned while trying to get check credentials: %v", err)
		apiErrors.APIErrorResponse(w, err)
		return
	}
	var token string
	token, err = jwtauth.NewToken(name)
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
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	log.Println("successfully returned token as cookie")
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := models.LogInRes{
		APIKey: apiKey,
	}
	json.NewEncoder(w).Encode(res)
	return
}
