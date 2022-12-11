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
	r.addRouterMiddleware(l)
	addAuthRoutes(api, a, l)
	addUserRoutes(api, u, l)
	addSearchRoutes(api, s, l)
	addBookmarkRoutes(api, b, l)

	return r
}

func (r *Router) initRouter() *mux.Router {
	api := r.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }).Methods("GET")
	return api
}

func (r *Router) addRouterMiddleware(l logs.Logger) {
	logmw := middleware.RouteLogger(l)
	r.router.Use(logmw)
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
	user.HandleFunc("/{APIKey}", auth.Authorized(handlers.DelUser(u, l), l)).Methods("DELETE")
	user.HandleFunc("/cmd/{APIKey}", auth.Authorized(handlers.GetCmds(u, l), l)).Methods("GET")
	user.HandleFunc("/cmd/{APIKey}", auth.Authorized(handlers.AddCmd(u, l), l)).Methods("POST")
	user.HandleFunc("/cmd/{APIKey}", auth.Authorized(handlers.DeleteCmd(u, l), l)).Methods("PATCH")
}

func addSearchRoutes(router *mux.Router, s search.Service, l logs.Logger) {
	search := router.PathPrefix("/search").Subrouter()
	search.HandleFunc("/{APIKey}/{args}", handlers.Search(s, l)).Methods("GET")
}

func addBookmarkRoutes(router *mux.Router, b bookmarks.Service, l logs.Logger) {
	bookmarks := router.PathPrefix("/bookmark").Subrouter()
	bookmarks.HandleFunc("/{APIKey}", auth.Authorized(handlers.GetAllBookmarks(b, l), l)).Methods("GET")
	bookmarks.HandleFunc("/{path}/{APIKey}", auth.Authorized(handlers.GetBookmarksFolder(b, l), l)).Methods("GET")
	bookmarks.HandleFunc("/{APIKey}", auth.Authorized(handlers.AddBookmark(b, l), l)).Methods("POST")
	bookmarks.HandleFunc("/file/{APIKey}", auth.Authorized(handlers.AddBookmarksFile(b, l), l)).Methods("POST")
	bookmarks.HandleFunc("/{APIKey}", auth.Authorized(handlers.DeleteBookmark(b, l), l)).Methods("DELETE")
}
