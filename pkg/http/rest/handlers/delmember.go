package handlers

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/conalli/bookshelf-backend/pkg/accounts"
// 	"github.com/conalli/bookshelf-backend/pkg/errors"
// )

// type delMemberResponse struct {
// 	TeamID            string `json:"teamId"`
// 	NumMembersDeleted int    `json:"numMembersDeleted"`
// }

// // DelMember is the handler for the delmember endpoint. Removes member from team.
// func DelMember(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("DelMember endpoint hit")
// 		var delMemberReq accounts.DelMemberRequest
// 		json.NewDecoder(r.Body).Decode(&delMemberReq)
// 		ok, err := t.DelMember(r.Context(), delMemberReq)
// 		if err != nil {
// 			log.Printf("error returned while trying to delete a member from team: %v", err)
// 			errors.APIErrorResponse(w, err)
// 			return
// 		}
// 		res := delSelfResponse{
// 			TeamID:            delMemberReq.TeamID,
// 			NumMembersDeleted: 0,
// 		}
// 		if !ok {
// 			log.Printf("failed to delete member from team: %s\n", delMemberReq.MemberName)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusBadRequest)
// 			json.NewEncoder(w).Encode(res)
// 			return
// 		}
// 		log.Printf("successfully deleted member from team: %s\n", delMemberReq.MemberName)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		res.NumMembersDeleted = 1
// 		json.NewEncoder(w).Encode(res)
// 		return
// 	}
// }

// // DelSelf is the handler for the delSelf endpoint. Removes member from team.
// func DelSelf(t accounts.TeamService) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("DelSelf endpoint hit")
// 		var delSelfReq accounts.DelSelfRequest
// 		json.NewDecoder(r.Body).Decode(&delSelfReq)
// 		ok, err := t.DelSelf(r.Context(), delSelfReq)
// 		if err != nil {
// 			log.Printf("error returned while trying to delete self from team: %v", err)
// 			errors.APIErrorResponse(w, err)
// 			return
// 		}
// 		res := delSelfResponse{
// 			TeamID:            delSelfReq.TeamID,
// 			NumMembersDeleted: 0,
// 		}
// 		if !ok {
// 			log.Printf("failed to delete self from team: %s\n", delSelfReq.ID)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusBadRequest)
// 			json.NewEncoder(w).Encode(res)
// 			return
// 		}
// 		log.Printf("successfully deleted self from team: %s\n", delSelfReq.ID)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		res.NumMembersDeleted = 1
// 		json.NewEncoder(w).Encode(res)
// 		return
// 	}
// }
