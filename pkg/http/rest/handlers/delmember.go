package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

type delMemberResponse struct {
	TeamID            string `json:"teamId"`
	NumMembersDeleted int    `json:"numMembersDeleted"`
}

// DelMember is the handler for the delmember endpoint. Removes member from team.
func DelMember(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DelMember endpoint hit")
		var delMemberReq request.DeleteMember
		json.NewDecoder(r.Body).Decode(&delMemberReq)
		ok, err := t.DeleteMember(r.Context(), delMemberReq)
		if err != nil {
			log.Printf("error returned while trying to delete a member from team: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		res := delMemberResponse{
			TeamID:            delMemberReq.TeamID,
			NumMembersDeleted: 0,
		}
		if !ok {
			log.Printf("failed to delete member from team: %s\n", delMemberReq.MemberName)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		log.Printf("successfully deleted member from team: %s\n", delMemberReq.MemberName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res.NumMembersDeleted = 1
		json.NewEncoder(w).Encode(res)
	}
}

// DelSelf is the handler for the delSelf endpoint. Removes member from team.
func DelSelf(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DelSelf endpoint hit")
		var delSelfReq request.DeleteSelf
		json.NewDecoder(r.Body).Decode(&delSelfReq)
		ok, err := t.DeleteSelf(r.Context(), delSelfReq)
		if err != nil {
			log.Printf("error returned while trying to delete self from team: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		res := delMemberResponse{
			TeamID:            delSelfReq.TeamID,
			NumMembersDeleted: 0,
		}
		if !ok {
			log.Printf("failed to delete self from team: %s\n", delSelfReq.ID)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		log.Printf("successfully deleted self from team: %s\n", delSelfReq.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res.NumMembersDeleted = 1
		json.NewEncoder(w).Encode(res)
	}
}
