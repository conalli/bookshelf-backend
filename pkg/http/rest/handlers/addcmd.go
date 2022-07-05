package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/gorilla/mux"
)

// AddCmdResponse represents the data returned upon successfully adding a cmd.
type AddCmdResponse struct {
	CmdsSet int    `json:"cmdsSet"`
	Cmd     string `json:"cmd"`
	URL     string `json:"url"`
}

// AddCmd is the handler for the setcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func AddCmd(u accounts.UserService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("SetCmd endpoint hit")
		vars := mux.Vars(r)
		APIKey := vars["APIKey"]
		var setCmdReq accounts.AddCmdRequest
		json.NewDecoder(r.Body).Decode(&setCmdReq)

		numUpdated, err := u.AddCmd(r.Context(), setCmdReq, APIKey)
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
		res := AddCmdResponse{
			CmdsSet: numUpdated,
			Cmd:     setCmdReq.Cmd,
			URL:     setCmdReq.URL,
		}
		json.NewEncoder(w).Encode(res)
	}
}
