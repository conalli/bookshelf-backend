package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/gorilla/mux"
)

// GetCmds is the handler for the getcmds endpoint. Checks credentials + JWT and if
// authorized returns all users cmds.
func GetCmds(u accounts.UserService, log logs.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		cmds, err := u.GetAllCmds(r.Context(), APIKey)
		if err != nil {
			log.Errorf("error returned while trying to get cmds: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Infof("successfully retrieved cmds: %v", cmds)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cmds)
	}
}
