package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
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
			errors.APIErrorResponse(w, errors.NewBadRequestError("invalid request url"))
			return
		}
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}

		user, apierr := a.OAuthRedirect(r.Context(), authProvider, authType, r.FormValue("code"), r.FormValue("state"), r.Cookies())
		if apierr != nil {
			log.Errorf("error returned while trying to create a new oauth user: %v", err)
			errors.APIErrorResponse(w, apierr)
			return
		}
		log.Infof("successfully created a new user: %+v", user)
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
