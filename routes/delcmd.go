package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/responses"
	"github.com/conalli/bookshelf-backend/models/user"
	"github.com/gorilla/mux"
)

// DelCmd is the handler for the delcmd endpoint. Checks credentials + JWT and if
// authorized deletes given cmd.
func DelCmd(w http.ResponseWriter, r *http.Request) {
	log.Println("DelCmd endpoint hit")
	vars := mux.Vars(r)
	apiKey := vars["apiKey"]
	var delCmdReq requests.DelCmdRequest
	json.NewDecoder(r.Body).Decode(&delCmdReq)

	result, err := user.DelCmd(r.Context(), delCmdReq, apiKey)
	if err != nil {
		log.Printf("error returned while trying to remove a cmd: %v", err)
		errors.APIErrorResponse(w, err)
		return
	}
	if result == 0 {
		log.Printf("could not remove cmd... maybe %s doesn't exists?", delCmdReq.Cmd)
		err := errors.NewBadRequestError("error: could not remove cmd")
		errors.APIErrorResponse(w, err)
		return
	}
	log.Printf("successfully updates cmds: %s, removed %d", delCmdReq.Cmd, result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := responses.DelCmdResponse{
		NumDeleted: result,
		Cmd:        delCmdReq.Cmd,
	}
	json.NewEncoder(w).Encode(res)
	return
}
