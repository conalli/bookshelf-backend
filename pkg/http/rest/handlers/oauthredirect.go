package handlers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

func OAuthRedirect(a auth.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
		}
		log.Infof("%+v", r.Form)
		a.OAuthFlow(r.Context(), r.FormValue("code"))
		w.WriteHeader(http.StatusOK)
	}
}
