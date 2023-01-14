package handlers

import (
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/gorilla/mux"
)

func OAuthRedirect(a auth.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		route := mux.Vars(r)
		authProvider, ok := route["authProvider"]
		authType, ok2 := route["authType"]
		if !(ok && ok2) {
			log.Error("no authType returned from redirect")
			apierr.APIErrorResponse(w, apierr.NewBadRequestError("invalid request url"))
			return
		}
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			apierr.APIErrorResponse(w, apierr.NewInternalServerError())
			return
		}
		tokens, apiErr := a.OAuthRedirect(r.Context(), authProvider, authType, r.FormValue("code"), r.FormValue("state"), r.Cookies())
		if apiErr != nil {
			log.Errorf("error returned while trying to %s a new oauth user: %v", authType, apiErr)
			apierr.APIErrorResponse(w, apiErr)
			return
		}
		cookies := tokens.NewTokenCookies(log, http.SameSiteLaxMode)
		log.Info("successfully returned token as cookie")
		auth.AddCookiesToResponse(w, cookies)
		http.Redirect(w, r, os.Getenv("ALLOWED_URL_DASHBOARD"), http.StatusTemporaryRedirect)
	}
}
