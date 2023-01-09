package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
)

// AddBookmarksFile attempts to add bookmarks to user from a given HTML file.
func AddBookmarksFile(b bookmarks.Service, log logs.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, ok := request.GetAPIKeyFromContext(r)
		if len(APIKey) < 1 || !ok {
			log.Error("could not get APIKey from context")
			errors.APIErrorResponse(w, errors.NewInternalServerError())
			return
		}
		if r.ContentLength > bookmarks.BookmarksFileMaxSize {
			log.Errorf("bookmarks file too large: %d, max: %d", r.ContentLength, bookmarks.BookmarksFileMaxSize)
			apierr := errors.NewAPIError(http.StatusExpectationFailed, "request too large", "bookmarks file too large")
			errors.APIErrorResponse(w, apierr)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, bookmarks.BookmarksFileMaxSize)
		err := r.ParseMultipartForm(200_000)
		if err != nil {
			log.Errorf("Could not parse multipart form: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
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
		res := struct {
			NumAdded int `json:"num_added"`
		}{
			NumAdded: num,
		}
		json.NewEncoder(w).Encode(res)
	}
}
