package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/team"
)

type delSelfResponse struct {
	TeamID            string `json:"teamId"`
	NumMembersDeleted int    `json:"numMembersDeleted"`
}

// DelSelf is the handler for the delSelf endpoint. Removes member from team.
func DelSelf(t team.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DelSelf endpoint hit")
		var delSelfReq team.DelSelfRequest
		json.NewDecoder(r.Body).Decode(&delSelfReq)
		ok, err := t.DelSelf(r.Context(), delSelfReq)
		if err != nil {
			log.Printf("error returned while trying to delete self from team: %v", err)
			errors.APIErrorResponse(w, err)
			return
		}
		res := delSelfResponse{
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
		return
	}
}
