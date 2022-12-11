package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/gorilla/mux"
)

// AddBookmarksFile attempts to add bookmarks to user from a given HTML file.
func AddBookmarksFile(b bookmarks.Service, log logs.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		err := r.ParseMultipartForm(200_000)
		if err != nil {
			log.Errorf("Could not parse multipart form: %v", err)
			w.WriteHeader(404)
			return
		}
		num, apierr := b.AddBookmarksFromFile(r.Context(), r, APIKey)
		if apierr != nil {
			log.Errorf("Could not add bookmarks from file: %v", apierr)
			errors.APIErrorResponse(w, apierr)
			return
		}
		log.Info(num)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := struct{ NumAdded int }{
			NumAdded: num,
		}
		json.NewEncoder(w).Encode(res)
	}
}
