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

// GetCmds is the handler for the getcmds endpoint. Checks credentials + JWT and if
// authorized returns all users cmds.
func GetCmds(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCmds endpoint hit")
	var getCmdsReq models.GetCmdsReq
	json.NewDecoder(r.Body).Decode(&getCmdsReq)

	if !jwtauth.Authorized(getCmdsReq.Name)(w, r) {
		jwtErr := apiErrors.NewApiError(http.StatusUnauthorized, apiErrors.ErrInvalidJWTToken.Error(), "error: invalid access token")
		apiErrors.APIErrorResponse(w, jwtErr)
		return
	}

	cmds, err := controllers.GetAllCmds(getCmdsReq)
	if err != nil {
		log.Printf("error returned while trying to get cmds: %v", err)
		apiErrors.APIErrorResponse(w, err)
		return
	}
	log.Printf("successfully retrieved cmds: %v", cmds)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cmds)
	return
}
