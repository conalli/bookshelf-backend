package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// CORSMiddleware sets CORS options for the main router.
func CORSMiddleware(r *mux.Router) http.Handler {
	allowedURL := []string{os.Getenv("ALLOWED_URL")}
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})
	origins := handlers.AllowedOrigins(allowedURL)
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	return handlers.CORS(methods, origins, headers)(r)
}
