package handlers

import (
	"net/http"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

func OAuthRequest(a auth.Service, log logs.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			apierr.APIErrorResponse(w, apierr.NewInternalServerError())
		}
		queryParams := r.URL.Query()
		authProvider := queryParams.Get("provider")
		authType := queryParams.Get("type")
		data, err := a.OAuthRequest(r.Context(), authProvider, authType)
		if err != nil {
			log.Error(err)
			apierr.APIErrorResponse(w, apierr.NewInternalServerError())
			return
		}
		stateCookie := &http.Cookie{
			Name:     "state",
			Value:    data.State,
			MaxAge:   120,
			HttpOnly: false,
			Secure:   true,
			Path:     "/api",
			Expires:  time.Now().Add(2 * time.Minute),
			SameSite: http.SameSiteNoneMode,
		}
		nonceCookie := &http.Cookie{
			Name:     "nonce",
			Value:    data.Nonce,
			MaxAge:   120,
			HttpOnly: false,
			Secure:   true,
			Path:     "/api",
			Expires:  time.Now().Add(2 * time.Minute),
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(w, stateCookie)
		http.SetCookie(w, nonceCookie)
		http.Redirect(w, r, data.AuthURL, http.StatusFound)
	}
}
