package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
)

// DeleteBookmarkResponse represents a successful response from the /user/bookmark POST endpoint.
type DeleteBookmarkResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	Path       string `json:"path,omitempty"`
	URL        string `json:"url"`
	NumDeleted int    `json:"num_deleted"`
}

// DeleteBookmark is the handler for the bookmark POST endpoint.
func DeleteBookmark(b bookmarks.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, ok := request.GetAPIKeyFromContext(r)
		if len(APIKey) < 1 || !ok {
			log.Error("could not get APIKey from context")
			apierr.APIErrorResponse(w, apierr.NewInternalServerError())
			return
		}
		delBookReq, parseErr := request.DecodeJSONRequest[request.DeleteBookmark](r.Body)
		if parseErr != nil {
			errRes := apierr.NewBadRequestError("could not parse request body")
			apierr.APIErrorResponse(w, errRes)
		}
		numUpdated, err := b.DeleteBookmark(r.Context(), delBookReq, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to delete a bookmark: %v", err)
			apierr.APIErrorResponse(w, err)
			return
		}
		if numUpdated == 0 {
			log.Error("could not delete bookmark")
			err := apierr.NewBadRequestError("error: could not add bookmark")
			apierr.APIErrorResponse(w, err)
			return
		}
		log.Info("successfully deleted bookmark")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := DeleteBookmarkResponse{
			ID:         delBookReq.ID,
			NumDeleted: numUpdated,
			Path:       delBookReq.Path,
			URL:        delBookReq.URL,
		}
		json.NewEncoder(w).Encode(res)
	}
}
