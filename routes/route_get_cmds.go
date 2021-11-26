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

func GetCmds(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCmds endpoint hit")
	if !jwtauth.Authorized()(w, r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var getCmdsReq models.GetCmdsReq
	json.NewDecoder(r.Body).Decode(&getCmdsReq)

	cmds, err := controllers.GetAllCmds(getCmdsReq)
	if err != nil {
		log.Printf("error returned while trying to get cmds: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status())
		getCmdsErr := apiErrors.ResError{
			Status: err.Status(),
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(getCmdsErr)
	} else {
		log.Printf("successfully retrieved cmds: %v", cmds)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cmds)
	}
}
