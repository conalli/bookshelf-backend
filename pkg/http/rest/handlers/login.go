package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

// LogIn is the handler for the login endpoint. Checks credentials and if
// correct returns JWT cookie.
func LogIn(a auth.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logInReq, parseErr := request.DecodeJSONRequest[request.LogIn](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		user, apierr := a.LogIn(r.Context(), logInReq)
		if apierr != nil {
			log.Errorf("error returned while trying to get check credentials: %v", apierr)
			errors.APIErrorResponse(w, apierr)
			return
		}
		tokens, err := auth.NewTokens(log, user.APIKey)
		if err != nil {
			log.Errorf("error returned while trying to create a new token: %v", err)
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		cookies := tokens.NewTokenCookies(log)
		log.Info("successfully returned token as cookie")
		auth.AddCookiesToResponse(w, cookies)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
