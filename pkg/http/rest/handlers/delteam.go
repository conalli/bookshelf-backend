package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

// DelTeam is the handler for the delteam endpoint. Checks for user role and if
// admin deletes team from database.
func DelTeam(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DelTeam endpoint hit")
		var delTeamReq accounts.DelTeamRequest
		json.NewDecoder(r.Body).Decode(&delTeamReq)
		numDeleted, err := t.Delete(r.Context(), delTeamReq)
		if err != nil {
			log.Printf("error returned while trying to delete a team: %v\n", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Printf("successfully deleted a team: %s\n", delTeamReq.TeamID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := struct {
			NumDeleted int `json:"numDeleted"`
		}{
			NumDeleted: numDeleted,
		}
		json.NewEncoder(w).Encode(res)
		return
	}
}
