package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/accounts"
	"github.com/conalli/bookshelf-backend/pkg/errors"
)

// NewTeam is the handler for the newteam endpoint. Checks db for team name and if
// unique adds new team with given data.
func NewTeam(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("NewTeam endpoint hit")
		var newTeamReq accounts.NewTeamRequest
		json.NewDecoder(r.Body).Decode(&newTeamReq)
		teamID, err := t.New(r.Context(), newTeamReq)
		if err != nil {
			log.Printf("error returned while trying to create a new team: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		log.Printf("successfully created a new team: %s", teamID)

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

// DelTeam is the handler for the delteam endpoint. Checks for user role and if
// admin deletes team from database.
func DelTeam(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DelTeam endpoint hit")
		var delTeamReq accounts.DelTeamRequest
		json.NewDecoder(r.Body).Decode(&delTeamReq)
		numDeleted, err := t.DeleteTeam(r.Context(), delTeamReq)
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
