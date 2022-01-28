package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/responses"
	"github.com/gorilla/mux"
)

// SetCmd is the handler for the setcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func SetCmd(w http.ResponseWriter, r *http.Request) {
	log.Println("SetCmd endpoint hit")
	vars := mux.Vars(r)
	user := vars["apiKey"]
	var setCmdReq requests.AddCmdRequest
	json.NewDecoder(r.Body).Decode(&setCmdReq)

	numUpdated, err := controllers.AddCmd(r.Context(), setCmdReq, user)
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
	res := responses.AddCmdResponse{
		CmdsSet: numUpdated,
		Cmd:     setCmdReq.Cmd,
		URL:     setCmdReq.URL,
	}
	json.NewEncoder(w).Encode(res)
	return
}
