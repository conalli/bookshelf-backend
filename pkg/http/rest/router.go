package rest

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/db/mongodb"
	"github.com/conalli/bookshelf-backend/pkg/http/middleware"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/search"
	"github.com/conalli/bookshelf-backend/pkg/team"
	"github.com/conalli/bookshelf-backend/pkg/user"
	"github.com/gorilla/mux"
)

// Router returns a router with all handlers assigned to it
func Router() *mux.Router {
	mongo := mongodb.New()
	u := user.NewService(mongo)
	s := search.NewService(mongo)
	t := team.NewService(mongo)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}).Methods("GET")
	router.HandleFunc("/user/signup", handlers.SignUp(u)).Methods("POST")
	router.HandleFunc("/user/login", handlers.LogIn(u)).Methods("POST")
	router.HandleFunc("/user/teams/{APIKey}", jwtauth.Authorized(handlers.GetAllTeams(u))).Methods("GET")
	router.HandleFunc("/user/cmds/{APIKey}", jwtauth.Authorized(handlers.GetCmds(u))).Methods("GET")
	router.HandleFunc("/user/addcmd/{APIKey}", jwtauth.Authorized(handlers.AddCmd(u))).Methods("PATCH")
	router.HandleFunc("/user/delcmd/{APIKey}", jwtauth.Authorized(handlers.DelCmd(u))).Methods("PATCH")
	router.HandleFunc("/user/del/{APIKey}", jwtauth.Authorized(handlers.DelUser(u))).Methods("DELETE")

	router.HandleFunc("/search/{APIKey}/{cmd}", handlers.Search(s)).Methods("GET")

	router.HandleFunc("/team/{APIKey}", jwtauth.Authorized(handlers.NewTeam(t))).Methods("POST")
	router.HandleFunc("/team/addmember/{APIKey}", jwtauth.Authorized(handlers.AddMember(t))).Methods("PATCH")
	router.HandleFunc("/team/delself/{APIKey}", jwtauth.Authorized(handlers.DelSelf(t))).Methods("PATCH")
	router.HandleFunc("/team/addcmd/{APIKey}", jwtauth.Authorized(handlers.AddTeamCmd(t))).Methods("PATCH")
	router.HandleFunc("/team/delcmd/{APIKey}", jwtauth.Authorized(handlers.DelTeamCmd(t))).Methods("PATCH")

	return router
}

// RouterWithCORS provides basic CORS middleware for a router.
func RouterWithCORS() http.Handler {
	router := Router()
	return middleware.CORSMiddleware(router)
}
