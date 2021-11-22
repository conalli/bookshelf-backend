package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
)

func GetCmds(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCmds endpoint hit")
	var getCmdsReq models.GetCmdsReq
	json.NewDecoder(r.Body).Decode(&getCmdsReq)

	cmds, err := controllers.GetAllCmds(getCmdsReq)
	if err != nil {
		log.Printf("error returned while trying to create a new user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		var testError tempError
		testError.Error = err.Error()
		json.NewEncoder(w).Encode(testError)
	} else {
		log.Printf("successfully retrieved cmds: %v", cmds)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cmds)
	}
}
