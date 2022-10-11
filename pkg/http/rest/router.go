package rest

import (
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
	log    logs.Logger
	router *mux.Router
}

// NewRouter returns a router with all handlers assigned to it
func NewRouter(l logs.Logger, v *validator.Validate, store db.Storage, cache db.Cache) *Router {
	u := accounts.NewUserService(l, v, store, cache)
	s := search.NewService(l, v, store, cache)

	r := &Router{l, mux.NewRouter()}
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
		r.log.Infof("Path:", tpl, "Err:", err1, "Methods:", met, "Err:", err2)
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
	user.HandleFunc("/cmd/{APIKey}", jwtauth.Authorized(handlers.GetCmds(u, l), l)).Methods("GET")
	user.HandleFunc("/cmd/{APIKey}", jwtauth.Authorized(handlers.AddCmd(u, l), l)).Methods("POST")
	user.HandleFunc("/cmd/{APIKey}", jwtauth.Authorized(handlers.DeleteCmd(u, l), l)).Methods("PATCH")
	user.HandleFunc("/bookmark/{APIKey}", jwtauth.Authorized(handlers.GetAllBookmarks(u, l), l)).Methods("GET")
	user.HandleFunc("/bookmark/{path}/{APIKey}", jwtauth.Authorized(handlers.GetBookmarksFolder(u, l), l)).Methods("GET")
	user.HandleFunc("/bookmark/{APIKey}", jwtauth.Authorized(handlers.AddBookmark(u, l), l)).Methods("POST")
	// user.HandleFunc("/bookmark/file/{APIKey}", jwtauth.Authorized(handlers.AddBookmarkFile(u, l), l)).Methods("POST")
	user.HandleFunc("/bookmark/{APIKey}", jwtauth.Authorized(handlers.DeleteBookmark(u, l), l)).Methods("DELETE")
}

func addSearchRoutes(router *mux.Router, s search.Service, l logs.Logger) {
	search := router.PathPrefix("/search").Subrouter()
	search.HandleFunc("/{APIKey}/{args}", handlers.Search(s, l)).Methods("GET")
}
