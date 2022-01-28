package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models/apiErrors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/responses"
	"github.com/gorilla/mux"
)

// DelUser is the handler for the delacc endpoint. Checks credentials and if
// authorized deletes user.
func DelUser(w http.ResponseWriter, r *http.Request) {
	log.Println("delacc endpoint hit")
	vars := mux.Vars(r)
	user := vars["apiKey"]
	var delAccReq requests.DelUserRequest
	json.NewDecoder(r.Body).Decode(&delAccReq)
	numDeleted, err := controllers.DelUser(r.Context(), delAccReq, user)
	if err != nil {
		log.Printf("error returned while trying to delete user: %v", err)
		apiErrors.APIErrorResponse(w, err)
		return
	}
	log.Printf("successfully deleted %d users: %v", numDeleted, delAccReq)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := responses.DelUserResponse{
		Name:       delAccReq.Name,
		NumDeleted: numDeleted,
	}
	json.NewEncoder(w).Encode(res)
	return
}
