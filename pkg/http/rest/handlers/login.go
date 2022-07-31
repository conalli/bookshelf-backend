package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

// LogInResponse represents the data returned upon successfully logging in.
type LogInResponse struct {
	ID     string `json:"id"`
	APIKey string `json:"APIKey"`
}

// LogIn is the handler for the login endpoint. Checks credentials and if
// correct returns JWT cookie for use with getcmds and setcmd.
func LogIn(u accounts.UserService, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Log In endpoint hit")
		logInReq, parseErr := request.DecodeJSONRequest[request.LogIn](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		currUser, err := u.LogIn(r.Context(), logInReq)
		if err != nil {
			log.Errorf("error returned while trying to get check credentials: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		tokens, err := jwtauth.NewTokens(currUser.APIKey, log)
		if err != nil {
			log.Errorf("error returned while trying to create a new token: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		accessToken := http.Cookie{
			Name:     "bookshelfjwt",
			Value:    tokens["access_token"],
			Path:     "/api",
			Expires:  time.Now().Add(15 * time.Minute),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		refreshToken := http.Cookie{
			Name:     "bookshelfrefresh",
			Value:    tokens["refresh_token"],
			Path:     "/api",
			Expires:  time.Now().Add(24 * time.Hour),
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		log.Info("successfully returned token as cookie")
		http.SetCookie(w, &accessToken)
		http.SetCookie(w, &refreshToken)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := LogInResponse{
			ID:     currUser.ID,
			APIKey: currUser.APIKey,
		}
		json.NewEncoder(w).Encode(res)
	}
}
