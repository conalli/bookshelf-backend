package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/team"
)

type addTeamCmdResponse struct {
	CmdsSet int    `json:"cmdsSet"`
	Cmd     string `json:"cmd"`
	URL     string `json:"url"`
}

// AddTeamCmd is the handler for the team/addcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func AddTeamCmd(t team.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("SetCmd endpoint hit")
		var setCmdReq team.AddTeamCmdRequest
		json.NewDecoder(r.Body).Decode(&setCmdReq)

		numUpdated, err := t.AddCmdToTeam(r.Context(), setCmdReq)
		if err != nil {
			log.Printf("error returned while trying to add a new cmd: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		if numUpdated == 0 {
			log.Printf("could not update cmds... maybe %s:%s already exists?", setCmdReq.Cmd, setCmdReq.URL)
			err := errors.NewBadRequestError("error: could not update cmds")
			errors.APIErrorResponse(w, err)
			return
		}
		log.Printf("successfully set cmd: %s, url: %s", setCmdReq.Cmd, setCmdReq.URL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := addTeamCmdResponse{
			CmdsSet: numUpdated,
			Cmd:     setCmdReq.Cmd,
			URL:     setCmdReq.URL,
		}
		json.NewEncoder(w).Encode(res)
		return
	}
}
