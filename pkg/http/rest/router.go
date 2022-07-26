package rest

import (
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/db"
	"github.com/conalli/bookshelf-backend/pkg/http/middleware"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/logs"
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
func NewRouter(l logs.Logger, v *validator.Validate, store db.Storage) *Router {
	u := accounts.NewUserService(l, v, store)
	s := search.NewService(l, v, store)

	r := &Router{mux.NewRouter()}
	api := r.initRouter()
	addUserRoutes(api, u, l)
	addSearchRoutes(api, s, l)

	return r
}

func (r *Router) initRouter() *mux.Router {
	api := r.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("Hello")) }).Methods("GET")
	return api
}

// Walk prints all the routes of the current router.
func (r *Router) Walk() *Router {
	r.router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		tpl, err1 := route.GetPathTemplate()
		met, err2 := route.GetMethods()
		log.Println("Path:", tpl, "Err:", err1, "Methods:", met, "Err:", err2)
		return nil
	})
	return r
}

// Handler returns the router as an http.Handler.
func (r *Router) Handler() http.Handler {
	return r.router
}

// HandlerWithCORS provides basic CORS middleware for a router.
func (r *Router) HandlerWithCORS() http.Handler {
	return middleware.CORSMiddleware(r.router)
}

func addUserRoutes(router *mux.Router, u accounts.UserService, l logs.Logger) {
	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("", handlers.SignUp(u, l)).Methods("POST")
	user.HandleFunc("/{APIKey}", jwtauth.Authorized(handlers.DelUser(u, l), l)).Methods("DELETE")
	user.HandleFunc("/login", handlers.LogIn(u, l)).Methods("POST")
	// user.HandleFunc("/teams/{APIKey}", jwtauth.Authorized(handlers.GetAllTeams(u))).Methods("GET")
	user.HandleFunc("/cmd/{APIKey}", jwtauth.Authorized(handlers.GetCmds(u, l), l)).Methods("GET")
	user.HandleFunc("/cmd/{APIKey}", jwtauth.Authorized(handlers.AddCmd(u, l), l)).Methods("POST")
	user.HandleFunc("/cmd/{APIKey}", jwtauth.Authorized(handlers.DeleteCmd(u, l), l)).Methods("PATCH")
	user.HandleFunc("/bookmark/{APIKey}", jwtauth.Authorized(handlers.GetAllBookmarks(u, l), l)).Methods("GET")
	user.HandleFunc("/bookmark/{path}/{APIKey}", jwtauth.Authorized(handlers.GetBookmarksFolder(u, l), l)).Methods("GET")
	user.HandleFunc("/bookmark/{APIKey}", jwtauth.Authorized(handlers.AddBookmark(u, l), l)).Methods("POST")
	user.HandleFunc("/bookmark/{APIKey}", jwtauth.Authorized(handlers.DeleteBookmark(u, l), l)).Methods("DELETE")
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

func addSearchRoutes(router *mux.Router, s search.Service, l logs.Logger) {
	search := router.PathPrefix("/search").Subrouter()
	search.HandleFunc("/{APIKey}/{cmd}", handlers.Search(s, l)).Methods("GET")
}
