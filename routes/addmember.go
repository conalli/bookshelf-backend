package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/controllers"
	"github.com/conalli/bookshelf-backend/models/errors"
	"github.com/conalli/bookshelf-backend/models/requests"
	"github.com/conalli/bookshelf-backend/models/responses"
)

// AddMember is the handler for the addmember endpoint. Checks db for team name and if
// unique adds new team with given data.
func AddMember(w http.ResponseWriter, r *http.Request) {
	log.Println("AddMember endpoint hit")
	var newMemberReq requests.AddMemberRequest
	json.NewDecoder(r.Body).Decode(&newMemberReq)
	ok, err := controllers.AddMember(r.Context(), newMemberReq)
	if err != nil {
		log.Printf("error returned while trying to create a new user: %v", err)
		errors.APIErrorResponse(w, err)
		return
	}
	res := responses.AddMemberResponse{
		TeamID:          newMemberReq.TeamID,
		NumMembersAdded: 0,
	}
	if !ok {
		log.Printf("failed to add member: %s\n", newMemberReq.MemberName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	log.Printf("successfully added a new member: %s\n", newMemberReq.MemberName)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res.NumMembersAdded = 1
	json.NewEncoder(w).Encode(res)
	return
}
