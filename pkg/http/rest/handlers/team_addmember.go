package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/accounts"
	"github.com/conalli/bookshelf-backend/pkg/errors"
)

type addMemberResponse struct {
	TeamID          string `json:"teamId"`
	NumMembersAdded int    `json:"numMembersAdded"`
}

// AddMember is the handler for the addmember endpoint. Checks db for team name and if
// unique adds new team with given data.
func AddMember(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("AddMember endpoint hit")
		var newMemberReq accounts.AddMemberRequest
		json.NewDecoder(r.Body).Decode(&newMemberReq)
		ok, err := t.AddMember(r.Context(), newMemberReq)
		if err != nil {
			log.Printf("error returned while trying to add a new member: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		res := addMemberResponse{
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
}
