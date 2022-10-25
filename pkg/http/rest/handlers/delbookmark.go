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

// DeleteBookmarkResponse represents a successful response from the /user/bookmark POST endpoint.
type DeleteBookmarkResponse struct {
	NumDeleted int    `json:"numDeleted"`
	Name       string `json:"name,omitempty"`
	Path       string `json:"path,omitempty"`
	URL        string `json:"url"`
}

// DeleteBookmark is the handler for the bookmark POST endpoint.
func DeleteBookmark(b bookmarks.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("DELETE BOOKMARK endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		delBookReq, parseErr := request.DecodeJSONRequest[request.DeleteBookmark](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		numUpdated, err := b.DeleteBookmark(r.Context(), delBookReq, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to delete a bookmark: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		if numUpdated == 0 {
			log.Error("could not delete bookmark")
			err := errors.NewBadRequestError("error: could not add bookmark")
			errors.APIErrorResponse(w, err)
			return
		}
		log.Info("successfully deleted bookmark")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := DeleteBookmarkResponse{
			NumDeleted: numUpdated,
			Path:       delBookReq.Path,
			URL:        delBookReq.URL,
		}
		json.NewEncoder(w).Encode(res)
	}
}
