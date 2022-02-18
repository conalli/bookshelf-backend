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

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/", handlers.SignUp(u)).Methods("POST")
	user.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.DelUser(u))).Methods("DELETE")
	user.HandleFunc("/login", handlers.LogIn(u)).Methods("POST")
	user.HandleFunc("/teams/{APIKey}", jwtauth.Authorized(handlers.GetAllTeams(u))).Methods("GET")
	user.HandleFunc("/cmds/{APIKey}", jwtauth.Authorized(handlers.GetCmds(u))).Methods("GET")
	user.HandleFunc("/addcmd/{APIKey}", jwtauth.Authorized(handlers.AddCmd(u))).Methods("PATCH")
	user.HandleFunc("/delcmd/{APIKey}", jwtauth.Authorized(handlers.DelCmd(u))).Methods("PATCH")

	team := router.PathPrefix("/team").Subrouter()
	team.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.NewTeam(t))).Methods("POST")
	team.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.DelTeam(t))).Methods("DELETE")
	team.HandleFunc("/addmember/{APIKey}", jwtauth.Authorized(handlers.AddMember(t))).Methods("PATCH")
	team.HandleFunc("/delself/{APIKey}", jwtauth.Authorized(handlers.DelSelf(t))).Methods("PATCH")
	team.HandleFunc("/delmember/{APIKey}", jwtauth.Authorized(handlers.DelMember(t))).Methods("PATCH")
	team.HandleFunc("/addcmd/{APIKey}", jwtauth.Authorized(handlers.AddTeamCmd(t))).Methods("PATCH")
	team.HandleFunc("/delcmd/{APIKey}", jwtauth.Authorized(handlers.DelTeamCmd(t))).Methods("PATCH")

	search := router.PathPrefix("/search").Subrouter()
	search.HandleFunc("/{APIKey}/{cmd}", handlers.Search(s)).Methods("GET")

	return router
}

// RouterWithCORS provides basic CORS middleware for a router.
func RouterWithCORS() http.Handler {
	router := Router()
	return middleware.CORSMiddleware(router)
}
