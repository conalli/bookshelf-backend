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
	router *mux.Router
}

// NewRouter returns a router with all handlers assigned to it
func NewRouter(v *validator.Validate, store db.Storage) *Router {
	u := accounts.NewUserService(v, store)
	s := search.NewService(store)

	r := &Router{mux.NewRouter()}
	api := r.initRouter()
	addUserRoutes(api, u)
	addSearchRoutes(api, s)

	return r
}

func (r *Router) initRouter() *mux.Router {
	api := r.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello")) }).Methods("GET")
	return api
}

// Walk prints all the routes of the current router.
func (r *Router) Walk() *Router {
	r.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err1 := route.GetPathTemplate()
		met, err2 := route.GetMethods()
		log.Println("Path:", tpl, "Err:", err1, "Methods:", met, "Err:", err2)
		return nil
	})
	return r
}

// Handler returns the router as an http.Handler.
func (r Router) Handler() http.Handler {
	return r.router
}

// HandlerWithCORS provides basic CORS middleware for a router.
func (r *Router) HandlerWithCORS() http.Handler {
	return middleware.CORSMiddleware(r.router)
}

func addUserRoutes(router *mux.Router, u accounts.UserService) {
	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("", handlers.SignUp(u)).Methods("POST")
	user.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.DelUser(u))).Methods("DELETE")
	user.HandleFunc("/login", handlers.LogIn(u)).Methods("POST")
	// user.HandleFunc("/teams/{APIKey}", jwtauth.Authorized(handlers.GetAllTeams(u))).Methods("GET")
	user.HandleFunc("/cmds/{APIKey}", jwtauth.Authorized(handlers.GetCmds(u))).Methods("GET")
	user.HandleFunc("/addcmd/{APIKey}", jwtauth.Authorized(handlers.AddCmd(u))).Methods("PATCH")
	user.HandleFunc("/delcmd/{APIKey}", jwtauth.Authorized(handlers.DeleteCmd(u))).Methods("PATCH")
}

// func addTeamRoutes(router *mux.Router, t accounts.TeamService) {
// 	team := router.PathPrefix("/team").Subrouter()
// 	team.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.NewTeam(t))).Methods("POST")
// 	team.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.DelTeam(t))).Methods("DELETE")
// 	team.HandleFunc("/addmember/{APIKey}", jwtauth.Authorized(handlers.AddMember(t))).Methods("PATCH")
// 	team.HandleFunc("/delself/{APIKey}", jwtauth.Authorized(handlers.DelSelf(t))).Methods("PATCH")
// 	team.HandleFunc("/delmember/{APIKey}", jwtauth.Authorized(handlers.DelMember(t))).Methods("PATCH")
// 	team.HandleFunc("/addcmd/{APIKey}", jwtauth.Authorized(handlers.AddTeamCmd(t))).Methods("PATCH")
// 	team.HandleFunc("/delcmd/{APIKey}", jwtauth.Authorized(handlers.DelTeamCmd(t))).Methods("PATCH")
// }

func addSearchRoutes(router *mux.Router, s search.Service) {
	search := router.PathPrefix("/search").Subrouter()
	search.HandleFunc("/{APIKey}/{cmd}", handlers.Search(s)).Methods("GET")
}
