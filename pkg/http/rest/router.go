package rest

import (
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/db"
	"github.com/conalli/bookshelf-backend/pkg/http/middleware"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/search"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Router wraps the *mux.Router type.
type Router struct {
	Router *mux.Router
}

// NewRouter returns a router with all handlers assigned to it
func NewRouter(v *validator.Validate, store db.Storage) *Router {
	u := accounts.NewUserService(v, store)
	// t := accounts.NewTeamService(repo)
	s := search.NewService(store)

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	}).Methods("GET")

	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("", handlers.SignUp(u)).Methods("POST")
	user.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.DelUser(u))).Methods("DELETE")
	user.HandleFunc("/login", handlers.LogIn(u)).Methods("POST")
	// user.HandleFunc("/teams/{APIKey}", jwtauth.Authorized(handlers.GetAllTeams(u))).Methods("GET")
	user.HandleFunc("/cmds/{APIKey}", jwtauth.Authorized(handlers.GetCmds(u))).Methods("GET")
	user.HandleFunc("/addcmd/{APIKey}", jwtauth.Authorized(handlers.AddCmd(u))).Methods("PATCH")
	user.HandleFunc("/delcmd/{APIKey}", jwtauth.Authorized(handlers.DeleteCmd(u))).Methods("PATCH")

	// team := router.PathPrefix("/team").Subrouter()
	// team.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.NewTeam(t))).Methods("POST")
	// team.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.DelTeam(t))).Methods("DELETE")
	// team.HandleFunc("/addmember/{APIKey}", jwtauth.Authorized(handlers.AddMember(t))).Methods("PATCH")
	// team.HandleFunc("/delself/{APIKey}", jwtauth.Authorized(handlers.DelSelf(t))).Methods("PATCH")
	// team.HandleFunc("/delmember/{APIKey}", jwtauth.Authorized(handlers.DelMember(t))).Methods("PATCH")
	// team.HandleFunc("/addcmd/{APIKey}", jwtauth.Authorized(handlers.AddTeamCmd(t))).Methods("PATCH")
	// team.HandleFunc("/delcmd/{APIKey}", jwtauth.Authorized(handlers.DelTeamCmd(t))).Methods("PATCH")

	search := api.PathPrefix("/search").Subrouter()
	search.HandleFunc("/{APIKey}/{cmd}", handlers.Search(s)).Methods("GET")

	return &Router{r}
}

// Walk prints all the routes of the current router.
func (r *Router) Walk() *Router {
	r.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err1 := route.GetPathTemplate()
		met, err2 := route.GetMethods()
		log.Println("Path:", tpl, "Err:", err1, "Methods:", met, "Err:", err2)
		return nil
	})
	return r
}

// WithCORS provides basic CORS middleware for a router.
func (r *Router) WithCORS() http.Handler {
	return middleware.CORSMiddleware(r.Router)
}
