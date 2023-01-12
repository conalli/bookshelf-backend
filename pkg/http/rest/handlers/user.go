package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

func GetUser(u accounts.UserService, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, ok := request.GetAPIKeyFromContext(r.Context())
		if !ok {
			log.Error("could not get APIKey from context")
			apierr.NewBadRequestError("could not get APIKey from auth token")
			return
		}
		user, apiErr := u.UserInfo(r.Context(), APIKey)
		if apiErr != nil {
			log.Error("could not get user info")
			apierr.APIErrorResponse(w, apiErr)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
