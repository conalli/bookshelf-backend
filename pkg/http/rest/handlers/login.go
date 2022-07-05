package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

// LogInResponse represents the data returned upon successfully logging in.
type LogInResponse struct {
	ID     string `json:"id"`
	APIKey string `json:"APIKey"`
}

// LogIn is the handler for the login endpoint. Checks credentials and if
// correct returns JWT cookie for use with getcmds and setcmd.
func LogIn(u accounts.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("LogIn endpoint hit")
		var logInReq accounts.LogInRequest
		json.NewDecoder(r.Body).Decode(&logInReq)

		currUser, err := u.LogIn(r.Context(), logInReq)
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
			Expires:  time.Now().Add(30 * time.Minute),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		log.Println("successfully returned token as cookie")
		http.SetCookie(w, &cookie)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := LogInResponse{
			ID:     currUser.ID,
			APIKey: currUser.APIKey,
		}
		json.NewEncoder(w).Encode(res)
	}
}
