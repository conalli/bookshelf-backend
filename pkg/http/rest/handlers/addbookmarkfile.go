package handlers

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
)

// AddBookmarkFile attempts to add bookmarks to user from a given HTML file.
func AddBookmarkFile(b bookmarks.Service, log logs.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("ADD BOOKMARK FILE endpoint hit")
		err := r.ParseMultipartForm(200_000)
		if err != nil {
			log.Errorf("Could not parse multipart form: %v", err)
			w.WriteHeader(404)
			return
		}
		num, err := b.AddBookmarksFromFile(r.Context(), r)
		if err != nil {
			return
		}
		log.Info(num)
		return
	}
}
