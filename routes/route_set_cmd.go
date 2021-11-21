package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models"
)

func SetCmd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SetCmd endpoint hit")
	var getCmdsReq models.SetCmdReq
	// Add error handling
	json.NewDecoder(r.Body).Decode(&getCmdsReq)
	numUpdated, err := controllers.AddCmd(getCmdsReq)
	if err != nil || numUpdated == 0 {
		fmt.Printf("error returned while trying to add a new cmd: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		var testError tempError
		testError.Error = err.Error()
		json.NewEncoder(w).Encode(testError)
	} else {
		fmt.Printf("successfully set cmd: %s, url: %s", getCmdsReq.Cmd, getCmdsReq.URL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := models.SetCmdRes{
			CmdsSet: numUpdated,
		}
		json.NewEncoder(w).Encode(res)
	}
}
