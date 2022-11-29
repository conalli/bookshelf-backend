package handlers

import (
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

func OAuthRequest(a auth.Service, log logs.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			errors.APIErrorResponse(w, errors.NewInternalServerError())
		}
		authType := r.FormValue("type")
		data, err := a.OAuthRequest(r.Context(), authType)
		if err != nil {
			log.Error(err)
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		stateCookie := &http.Cookie{
			Name:     "state",
			Value:    data.State,
			MaxAge:   120,
			HttpOnly: true,
			Path:     "/api",
			Expires:  time.Now().Add(2 * time.Minute),
		}
		nonceCookie := &http.Cookie{
			Name:     "nonce",
			Value:    data.Nonce,
			MaxAge:   120,
			HttpOnly: true,
			Path:     "/api",
			Expires:  time.Now().Add(2 * time.Minute),
		}
		http.SetCookie(w, stateCookie)
		http.SetCookie(w, nonceCookie)
		http.Redirect(w, r, data.AuthURL, http.StatusFound)
	}
}
