package handlers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

func OAuthRedirect(a auth.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		stateCookie, err := r.Cookie("state")
		if err != nil {
			log.Error(err)
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		nonceCookie, err := r.Cookie("nonce")
		if err != nil {
			log.Error(err)
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		_, apierr := a.OAuthRedirect(r.Context(), r.FormValue("code"), r.FormValue("state"), stateCookie, nonceCookie)
		if apierr != nil {
			log.Error(err)
			errors.APIErrorResponse(w, apierr)
		}
		w.WriteHeader(http.StatusOK)
	}
}
