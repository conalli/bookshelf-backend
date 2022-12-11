package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
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
		APIKey := vars["APIKey"]
		books, err := b.GetBookmarksFolder(r.Context(), path, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to get cmds: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Infof("successfully retrieved cmds: %v", books)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(books)
	}
}
