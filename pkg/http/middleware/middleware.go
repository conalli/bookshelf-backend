package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
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

// Authorized reads the JWT from the incoming request and returns whether the user is authorized or not.
func Authorized(log logs.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookies := r.Cookies()
			if len(cookies) < 1 {
				log.Error("no cookies in request")
				errors.APIErrorResponse(w, errors.NewBadRequestError("no cookies in request"))
				return
			}
			bookshelfCookies, err := request.FindCookies(cookies, auth.BookshelfTokenCode, auth.BookshelfAccessToken, auth.BookshelfRefreshToken)
			if err != nil {
				log.Errorf("could not find bookshelf cookies: %v", err)
				errors.APIErrorResponse(w, errors.NewBadRequestError("could not find bookshelf cookies"))
				return
			}
			log.Info(bookshelfCookies)
			code := bookshelfCookies[auth.BookshelfAccessToken].Value
			accessToken := bookshelfCookies[auth.BookshelfAccessToken].Value
			jwt, err := auth.ParseAccessToken(log, accessToken, code)
			if err != nil {
				log.Error("could not parse access token: %v", err)
				errors.APIErrorResponse(w, errors.NewJWTTokenError(err.Error()))
			}
			r.WithContext(context.WithValue(r.Context(), request.JWTAPIKey, jwt.Subject))
			next.ServeHTTP(w, r)
		})
	}
}
