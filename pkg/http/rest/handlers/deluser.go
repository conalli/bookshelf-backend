package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/gorilla/mux"
)

// DelUserResponse represents the data returned upon successfully deleting a user.
type DelUserResponse struct {
	Name       string `json:"name"`
	NumDeleted int    `json:"numDeleted"`
}

// DelUser is the handler for the delacc endpoint. Checks credentials and if
// authorized deletes user.
func DelUser(u accounts.UserService, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("DELETE USER endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		delAccReq, parseErr := request.DecodeJSONRequest[request.DeleteUser](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		numDeleted, err := u.Delete(r.Context(), delAccReq, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to delete user: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Infof("successfully deleted %d users: %v", numDeleted, delAccReq)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := DelUserResponse{
			Name:       delAccReq.Name,
			NumDeleted: numDeleted,
		}
		json.NewEncoder(w).Encode(res)
	}
}
