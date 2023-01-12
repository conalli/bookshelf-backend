package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/gorilla/mux"
)

// GetBookmarksFolder is the handler for the /user/bookmarks/{Path} GET endpoint. Checks credentials + JWT and if
// authorized returns all users bookmarks.
func GetBookmarksFolder(b bookmarks.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]
		APIKey, ok := request.GetAPIKeyFromContext(r.Context())
		if len(APIKey) < 1 || !ok {
			log.Error("could not get APIKey from context")
			apierr.APIErrorResponse(w, apierr.NewInternalServerError())
			return
		}
		books, err := b.GetBookmarksFolder(r.Context(), path, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to get cmds: %v", err)
			apierr.APIErrorResponse(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(books)
	}
}
