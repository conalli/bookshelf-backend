package handlers

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/conalli/bookshelf-backend/pkg/accounts"
// 	"github.com/conalli/bookshelf-backend/pkg/errors"
// 	"github.com/gorilla/mux"
// )

// // GetAllTeams is the handler for the getteams endpoint. Checks credentials + JWT and if
// // authorized returns all teams.
// func GetAllTeams(u accounts.UserService) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("getteams endpoint hit")
// 		vars := mux.Vars(r)
// 		APIKey := vars["APIKey"]
// 		teams, err := u.GetTeams(r.Context(), APIKey)
// 		if err != nil {
// 			log.Printf("error returned while trying to get teams: %v", err)
// 			errors.APIErrorResponse(w, err)
// 			return
// 		}
// 		log.Printf("successfully retrieved teams: %v", teams)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(teams)
// 		return
// 	}
// }
