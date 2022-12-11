package rest

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/db"
	"github.com/conalli/bookshelf-backend/pkg/http/middleware"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/conalli/bookshelf-backend/pkg/services/search"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Router wraps the *mux.Router type.
type Router struct {
	log    logs.Logger
	router *mux.Router
}

// NewRouter returns a router with all handlers assigned to it
func NewRouter(l logs.Logger, v *validator.Validate, store db.Storage, cache db.Cache, p *oidc.Provider) *Router {
	a := auth.NewService(l, v, p, store)
	u := accounts.NewUserService(l, v, store, cache)
	s := search.NewService(l, v, store, cache)
	b := bookmarks.NewService(l, v, store)
	r := &Router{l, mux.NewRouter()}

	api := r.initRouter()
	addAuthRoutes(api, a, l)
	addUserRoutes(api, u, l)
	addSearchRoutes(api, s, l)
	addBookmarkRoutes(api, b, l)

	r.router.Use(middleware.RouteLogger(l))
	return r
}

func (r *Router) initRouter() *mux.Router {
	api := r.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }).Methods("GET")
	return api
}

// Walk prints all the routes of the current router.
func (r *Router) Walk() *Router {
	r.router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		tpl, err1 := route.GetPathTemplate()
		met, err2 := route.GetMethods()
		r.log.Info("Path:", tpl, "Err:", err1, "Methods:", met, "Err:", err2)
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
	return middleware.CORS(r.router)
}

func addAuthRoutes(router *mux.Router, a auth.Service, l logs.Logger) {
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", handlers.SignUp(a, l)).Methods("POST")
	auth.HandleFunc("/login", handlers.LogIn(a, l)).Methods("POST")
	auth.HandleFunc("/oauth", handlers.OAuthRequest(a, l)).Methods("GET")
	auth.HandleFunc("/redirect/{authProvider}/{authType}", handlers.OAuthRedirect(a, l)).Methods("GET")
}

func addUserRoutes(router *mux.Router, u accounts.UserService, l logs.Logger) {
	user := router.PathPrefix("/user").Subrouter()
	user.Use(middleware.Authorized(l))
	user.HandleFunc("/{APIKey}", handlers.DelUser(u, l)).Methods("DELETE")
	user.HandleFunc("/cmd/{APIKey}", handlers.GetCmds(u, l)).Methods("GET")
	user.HandleFunc("/cmd/{APIKey}", handlers.AddCmd(u, l)).Methods("POST")
	user.HandleFunc("/cmd/{APIKey}", handlers.DeleteCmd(u, l)).Methods("PATCH")
}

func addSearchRoutes(router *mux.Router, s search.Service, l logs.Logger) {
	search := router.PathPrefix("/search").Subrouter()
	search.HandleFunc("/{APIKey}/{args}", handlers.Search(s, l)).Methods("GET")
}

func addBookmarkRoutes(router *mux.Router, b bookmarks.Service, l logs.Logger) {
	bookmarks := router.PathPrefix("/bookmark").Subrouter()
	bookmarks.Use(middleware.Authorized(l))
	bookmarks.HandleFunc("/{APIKey}", handlers.GetAllBookmarks(b, l)).Methods("GET")
	bookmarks.HandleFunc("/{path}/{APIKey}", handlers.GetBookmarksFolder(b, l)).Methods("GET")
	bookmarks.HandleFunc("/{APIKey}", handlers.AddBookmark(b, l)).Methods("POST")
	bookmarks.HandleFunc("/file/{APIKey}", handlers.AddBookmarksFile(b, l)).Methods("POST")
	bookmarks.HandleFunc("/{APIKey}", handlers.DeleteBookmark(b, l)).Methods("DELETE")
}
