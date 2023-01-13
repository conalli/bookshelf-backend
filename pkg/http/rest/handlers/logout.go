package handlers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
)

func LogOut(a auth.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, _, ok := request.GetSearchKeysFromContext(r.Context())
		if !ok {
			log.Error("could not get APIKey from context")
			apierr.APIErrorResponse(w, apierr.NewInternalServerError())
			return
		}
		apiErr := a.LogOut(r.Context(), APIKey)
		if apiErr != nil {
			log.Errorf("error attempting to log out: %+v", apiErr)
			apierr.APIErrorResponse(w, apiErr)
			return
		}
		auth.RemoveBookshelfCookies(w)
		w.WriteHeader(http.StatusOK)
	}
}
