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
		newUserReq, parseErr := request.DecodeJSONRequest[request.SignUp](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		tokens, user, apierr := a.SignUp(r.Context(), newUserReq)
		if apierr != nil {
			log.Errorf("error returned while trying to create a new user: %v", apierr)
			errors.APIErrorResponse(w, apierr)
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
