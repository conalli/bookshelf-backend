package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
	"github.com/gorilla/mux"
)

// SetCmd is the handler for the setcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func SetCmd(w http.ResponseWriter, r *http.Request) {
	log.Println("SetCmd endpoint hit")
	vars := mux.Vars(r)
	user := vars["apiKey"]
	var setCmdReq models.SetCmdRequest
	json.NewDecoder(r.Body).Decode(&setCmdReq)

	numUpdated, err := controllers.AddCmd(r.Context(), setCmdReq, user)
	if err != nil {
		log.Printf("error returned while trying to add a new cmd: %v", err)
		apiErrors.APIErrorResponse(w, err)
		return
	}
	if numUpdated == 0 {
		log.Printf("could not update cmds... maybe %s:%s already exists?", setCmdReq.Cmd, setCmdReq.URL)
		err := apiErrors.NewBadRequestError("error: could not update cmds")
		apiErrors.APIErrorResponse(w, err)
		return
	}
	log.Printf("successfully set cmd: %s, url: %s", setCmdReq.Cmd, setCmdReq.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := models.SetCmdResponse{
		CmdsSet: numUpdated,
		Cmd:     setCmdReq.Cmd,
		URL:     setCmdReq.URL,
	}
	json.NewEncoder(w).Encode(res)
	return
}
