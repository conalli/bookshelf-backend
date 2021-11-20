package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
)

func GetCmds(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetCmds endpoint hit")
	var getCmdsReq models.GetCmdsReq
	json.NewDecoder(r.Body).Decode(&getCmdsReq)

	cmds, err := controllers.GetAllCmds(getCmdsReq)
	if err != nil {
		fmt.Printf("error returned while trying to create a new user: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		var testError tempError
		testError.Error = err.Error()
		json.NewEncoder(w).Encode(testError)
	} else {
		fmt.Printf("successfully retrieved cmds: %v", cmds)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cmds)
	}
}
