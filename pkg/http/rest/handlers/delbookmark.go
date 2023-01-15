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

// DeleteBookmarkResponse represents a successful response from the /user/bookmark POST endpoint.
type DeleteBookmarkResponse struct {
	ID         string `json:"id"`
	NumDeleted int    `json:"num_deleted"`
}

// DeleteBookmark is the handler for the bookmark POST endpoint.
func DeleteBookmark(b bookmarks.Service, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, ok := request.GetAPIKeyFromContext(r.Context())
		if len(APIKey) < 1 || !ok {
			log.Error("could not get APIKey from context")
			apierr.APIErrorResponse(w, apierr.NewInternalServerError())
			return
		}
		bookmarkID := mux.Vars(r)["id"]
		numUpdated, err := b.DeleteBookmark(r.Context(), bookmarkID, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to delete a bookmark: %v", err)
			apierr.APIErrorResponse(w, err)
			return
		}
		if numUpdated == 0 {
			log.Error("could not delete bookmark")
			err := apierr.NewBadRequestError("error: could not delete bookmark")
			apierr.APIErrorResponse(w, err)
			return
		}
		log.Info("successfully deleted bookmark")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := DeleteBookmarkResponse{
			ID:         bookmarkID,
			NumDeleted: numUpdated,
		}
		json.NewEncoder(w).Encode(res)
	}
}
