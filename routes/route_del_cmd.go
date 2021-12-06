package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
)

// DelCmd is the handler for the delcmd endpoint. Checks credentials + JWT and if
// authorized deletes given cmd.
func DelCmd(w http.ResponseWriter, r *http.Request) {
	log.Println("DelCmd endpoint hit")
	var delCmdReq models.DelCmdReq
	json.NewDecoder(r.Body).Decode(&delCmdReq)

	result, err := controllers.DelCmd(r.Context(), delCmdReq)
	if err != nil {
		log.Printf("error returned while trying to remove a cmd: %v", err)
		apiErrors.APIErrorResponse(w, err)
		return
	}
	log.Printf("successfully updates cmds: %s, removed %d", delCmdReq.Cmd, result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := models.DelCmdRes{
		NumDeleted: result,
		Cmd:        delCmdReq.Cmd,
	}
	json.NewEncoder(w).Encode(res)
	return
}
