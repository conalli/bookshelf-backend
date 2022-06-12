package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/accounts"
	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/gorilla/mux"
)

// GetCmds is the handler for the getcmds endpoint. Checks credentials + JWT and if
// authorized returns all users cmds.
func GetCmds(u accounts.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("getcmds endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		cmds, err := u.GetAllCmds(r.Context(), APIKey)
		if err != nil {
			log.Printf("error returned while trying to get cmds: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Printf("successfully retrieved cmds: %v", cmds)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cmds)
		return
	}
}
