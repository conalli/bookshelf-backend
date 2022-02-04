package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/team"
)

// NewTeam is the handler for the newteam endpoint. Checks db for team name and if
// unique adds new team with given data.
func NewTeam(t team.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("NewTeam endpoint hit")
		var newTeamReq team.NewTeamRequest
		json.NewDecoder(r.Body).Decode(&newTeamReq)
		teamID, err := t.New(r.Context(), newTeamReq)
		if err != nil {
			log.Printf("error returned while trying to create a new user: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Printf("successfully created a new user: %s", teamID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := struct {
			ID string `json:"id"`
		}{
			ID: teamID,
		}
		json.NewEncoder(w).Encode(res)
		return
	}
}
