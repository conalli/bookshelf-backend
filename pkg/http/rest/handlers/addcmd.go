package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// AddCmdResponse represents the data returned upon successfully adding a cmd.
type AddCmdResponse struct {
	CmdsSet int    `json:"cmdsSet"`
	Cmd     string `json:"cmd"`
	URL     string `json:"url"`
}

// AddCmd is the handler for the setcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func AddCmd(u accounts.UserService, log *zap.SugaredLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("ADD CMD endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		setCmdReq, parseErr := request.DecodeJSONRequest[request.AddCmd](r.Body)
		if parseErr != nil {
			errRes := errors.NewBadRequestError("could not parse request body")
			errors.APIErrorResponse(w, errRes)
		}
		numUpdated, err := u.AddCmd(r.Context(), setCmdReq, APIKey)
		if err != nil {
			log.Errorf("error returned while trying to add a new cmd: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		if numUpdated == 0 {
			log.Errorf("could not update cmds... maybe %s:%s already exists?", setCmdReq.Cmd, setCmdReq.URL)
			err := errors.NewBadRequestError("error: could not update cmds")
			errors.APIErrorResponse(w, err)
			return
		}
		log.Infof("successfully set cmd: %s, url: %s", setCmdReq.Cmd, setCmdReq.URL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := AddCmdResponse{
			CmdsSet: numUpdated,
			Cmd:     setCmdReq.Cmd,
			URL:     setCmdReq.URL,
		}
		json.NewEncoder(w).Encode(res)
	}
}

type addTeamCmdResponse struct {
	CmdsSet int    `json:"cmdsSet"`
	Cmd     string `json:"cmd"`
	URL     string `json:"url"`
}

// AddTeamCmd is the handler for the team/addcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func AddTeamCmd(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("SetCmd endpoint hit")
		var setCmdReq request.AddTeamCmd
		json.NewDecoder(r.Body).Decode(&setCmdReq)

		numUpdated, err := t.AddCmd(r.Context(), setCmdReq)
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
	}
}
