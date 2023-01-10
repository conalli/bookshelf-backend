package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
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
			apierr.APIErrorResponse(w, apierr.NewInternalServerError())
			return
		}
		if r.ContentLength > bookmarks.BookmarksFileMaxSize {
			log.Errorf("bookmarks file too large: %d, max: %d", r.ContentLength, bookmarks.BookmarksFileMaxSize)
			apiErr := apierr.NewAPIError(http.StatusExpectationFailed, errors.New("request too large"), "bookmarks file too large")
			apierr.APIErrorResponse(w, apiErr)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, bookmarks.BookmarksFileMaxSize)
		err := r.ParseMultipartForm(200_000)
		if err != nil {
			log.Errorf("Could not parse multipart form: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		num, apiErr := b.AddBookmarksFromFile(r.Context(), r, APIKey)
		if apiErr != nil {
			log.Errorf("Could not add bookmarks from file: %v", apiErr)
			apierr.APIErrorResponse(w, apiErr)
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
