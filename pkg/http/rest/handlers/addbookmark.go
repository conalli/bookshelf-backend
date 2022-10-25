package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/gorilla/mux"
)

// AddBookmarkResponse represents a successful response from the /user/bookmark POST endpoint.
type AddBookmarkResponse struct {
	NumAdded int    `json:"numAdded"`
	Name     string `json:"name,omitempty"`
	Path     string `json:"path"`
	URL      string `json:"url"`
}

// AddBookmark is the handler for the bookmark POST endpoint.
func AddBookmark(b bookmarks.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("ADD BOOKMARK endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		addBookReq, parseErr := request.DecodeJSONRequest[request.AddBookmark](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
			return
		}
		numUpdated, err := b.AddBookmark(r.Context(), addBookReq, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to add a new bookmark: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		if numUpdated == 0 {
			log.Error("could not add bookmark")
			err := errors.NewBadRequestError("error: could not add bookmark")
			errors.APIErrorResponse(w, err)
			return
		}
		log.Info("successfully added bookmark")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := AddBookmarkResponse{
			NumAdded: numUpdated,
			Name:     addBookReq.Name,
			Path:     addBookReq.Path,
			URL:      addBookReq.URL,
		}
		json.NewEncoder(w).Encode(res)
	}
}
