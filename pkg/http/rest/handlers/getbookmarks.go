package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/gorilla/mux"
)

// GetAllBookmarks is the handler for the /user/bookmarks GET endpoint. Checks credentials + JWT and if
// authorized returns all users bookmarks.
func GetAllBookmarks(u accounts.UserService, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("GET BOOKMARKS endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		books, err := u.GetAllBookmarks(r.Context(), APIKey)
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
