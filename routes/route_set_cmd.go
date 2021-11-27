package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
	"github.com/conalli/bookshelf-backend/utils/auth/jwtauth"
)

// SetCmd is the handler for the setcmd endpoint. Checks credentials + JWT and if
// authorized sets new cmd.
func SetCmd(w http.ResponseWriter, r *http.Request) {
	log.Println("SetCmd endpoint hit")
	var setCmdReq models.SetCmdReq
	json.NewDecoder(r.Body).Decode(&setCmdReq)

	if !jwtauth.Authorized(setCmdReq.Name)(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	numUpdated, err := controllers.AddCmd(setCmdReq)
	if err != nil || numUpdated == 0 {
		log.Printf("error returned while trying to add a new cmd: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status())
		setCmdError := apiErrors.ResError{
			Status: err.Status(),
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(setCmdError)
		return
	}
	log.Printf("successfully set cmd: %s, url: %s", setCmdReq.Cmd, setCmdReq.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := models.SetCmdRes{
		CmdsSet: numUpdated,
	}
	json.NewEncoder(w).Encode(res)
	return
}
