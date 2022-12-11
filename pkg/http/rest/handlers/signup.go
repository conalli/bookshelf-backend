package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

// SignUp is the handler for the signup endpoint. Checks db for username and if
// unique adds new user with given credentials.
func SignUp(a auth.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Sign Up endpoint hit")
		newUserReq, parseErr := request.DecodeJSONRequest[request.SignUp](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		user, apierr := a.SignUp(r.Context(), newUserReq)
		if apierr != nil {
			log.Errorf("error returned while trying to create a new user: %v", apierr)
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
		for _, cookie := range cookies {
			http.SetCookie(w, cookie)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
