package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"go.uber.org/zap"
)

// SignUp is the handler for the signup endpoint. Checks db for username and if
// unique adds new user with given credentials.
func SignUp(u accounts.UserService, log *zap.SugaredLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Sign Up endpoint hit")
		newUserReq, parseErr := request.DecodeJSONRequest[request.SignUp](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		newUser, err := u.NewUser(r.Context(), newUserReq)
		if err != nil {
			log.Errorf("error returned while trying to create a new user: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Infof("successfully created a new user: %+v", newUser)
		var token string
		token, err = jwtauth.NewToken(newUser.APIKey)
		if err != nil {
			log.Errorf("error returned while trying to create a new token: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		// Use Secure during production.
		cookie := http.Cookie{
			Name:     "bookshelfjwt",
			Value:    token,
			Expires:  time.Now().Add(30 * time.Minute),
			HttpOnly: true,
			// Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		log.Info("successfully returned token as cookie")
		http.SetCookie(w, &cookie)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := newUser
		json.NewEncoder(w).Encode(res)
	}
}

// NewTeam is the handler for the newteam endpoint. Checks db for team name and if
// unique adds new team with given data.
func NewTeam(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("NewTeam endpoint hit")
		var newTeamReq request.NewTeam
		json.NewDecoder(r.Body).Decode(&newTeamReq)
		teamID, err := t.New(r.Context(), newTeamReq)
		if err != nil {
			log.Printf("error returned while trying to create a new team: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Printf("successfully created a new team: %s", teamID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := struct {
			ID string `json:"id"`
		}{
			ID: teamID,
		}
		json.NewEncoder(w).Encode(res)
	}
}
