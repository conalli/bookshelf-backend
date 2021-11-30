package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// SetCmd is the handler for the setcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func SetCmd(w http.ResponseWriter, r *http.Request) {
	log.Println("SetCmd endpoint hit")
	var setCmdReq models.SetCmdReq
	json.NewDecoder(r.Body).Decode(&setCmdReq)

	numUpdated, err := controllers.AddCmd(setCmdReq)
	if err != nil || numUpdated == 0 {
		log.Printf("error returned while trying to add a new cmd: %v", err)
		apiErrors.APIErrorResponse(w, err)
		return
	}
	log.Printf("successfully set cmd: %s, url: %s", setCmdReq.Cmd, setCmdReq.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := models.SetCmdRes{
		CmdsSet: numUpdated,
		Cmd:     setCmdReq.Cmd,
		URL:     setCmdReq.URL,
	}
	json.NewEncoder(w).Encode(res)
	return
}
