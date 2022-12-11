package middleware

import (
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// CORS sets CORS options for the main router.
func CORS(h http.Handler) http.Handler {
	allowedURL := []string{os.Getenv("ALLOWED_URL_BASE"), os.Getenv("ALLOWED_URL_DASHBOARD")}
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE"})
	origins := handlers.AllowedOrigins(allowedURL)
	credentials := handlers.AllowCredentials()
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	return handlers.CORS(methods, origins, credentials, headers)(h)
}

func RouteLogger(log logs.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Infof("%s %s %s %s", r.RemoteAddr, r.Method, r.URL.String(), r.Proto)
			next.ServeHTTP(w, r)
		})
	}
}
