package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
	"github.com/conalli/bookshelf-backend/utils/apiErrors"
)

func SetCmd(w http.ResponseWriter, r *http.Request) {
	log.Println("SetCmd endpoint hit")
	var getCmdsReq models.SetCmdReq
	// Add error handling
	json.NewDecoder(r.Body).Decode(&getCmdsReq)
	numUpdated, err := controllers.AddCmd(getCmdsReq)
	if err != nil || numUpdated == 0 {
		log.Printf("error returned while trying to add a new cmd: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status())
		setCmdError := apiErrors.ResError{
			Status: err.Status(),
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(setCmdError)
	} else {
		log.Printf("successfully set cmd: %s, url: %s", getCmdsReq.Cmd, getCmdsReq.URL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := models.SetCmdRes{
			CmdsSet: numUpdated,
		}
		json.NewEncoder(w).Encode(res)
	}
}
