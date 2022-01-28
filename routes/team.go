package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
	"github.com/conalli/bookshelf-backend/models/requests"
)

// NewTeam is the handler for the newteam endpoint. Checks db for team name and if
// unique adds new team with given data.
func NewTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("NewTeam endpoint hit")
	var newTeamReq requests.NewTeamReq
	json.NewDecoder(r.Body).Decode(&newTeamReq)
	teamID, err := controllers.CreateNewTeam(r.Context(), newTeamReq)
	if err != nil {
		log.Printf("error returned while trying to create a new user: %v", err)
		apiErrors.APIErrorResponse(w, err)
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
