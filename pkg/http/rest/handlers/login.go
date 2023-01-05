package handlers

import (
	"net/http"
	"os"

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
		tokens, apierr := a.LogIn(r.Context(), logInReq)
		if apierr != nil {
			log.Errorf("error returned while trying to get check credentials: %v", apierr)
			errors.APIErrorResponse(w, apierr)
			return
		}
		cookies := tokens.NewTokenCookies(log)
		log.Info("successfully returned token as cookie")
		auth.AddCookiesToResponse(w, cookies)
		http.Redirect(w, r, os.Getenv("ALLOWED_URL_DASHBOARD"), http.StatusFound)
	}
}
