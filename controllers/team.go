package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/team"
)

// NewTeam is the handler for the newteam endpoint. Checks db for team name and if
// unique adds new team with given data.
func NewTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("NewTeam endpoint hit")
	var newTeamReq requests.NewTeamRequest
	json.NewDecoder(r.Body).Decode(&newTeamReq)
	teamID, err := team.New(r.Context(), newTeamReq)
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
