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
		logInReq, err := request.DecodeJSONRequest[request.LogIn](r.Body)
		if err != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		authUser, apierr := a.LogIn(r.Context(), logInReq)
		if apierr != nil {
			log.Errorf("error returned while trying to get check credentials: %v", apierr)
			errors.APIErrorResponse(w, apierr)
			return
		}
		cookies := authUser.Tokens.NewTokenCookies(log)
		log.Info("successfully returned token as cookie")
		auth.AddCookiesToResponse(w, cookies)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authUser.User)
	}
}
