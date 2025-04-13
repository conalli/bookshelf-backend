package middleware

import (
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
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
				apierr.APIErrorResponse(w, apierr.NewUnauthorizedError("no cookies in request"))
				return
			}
			bookshelfCookies, err := request.FindCookies(cookies, auth.BookshelfTokenCode, auth.BookshelfAccessToken)
			if err != nil {
				log.Errorf("could not find bookshelf cookies: %v", err)
				apierr.APIErrorResponse(w, apierr.NewUnauthorizedError("no cookies in request"))
				return
			}
			accessToken := bookshelfCookies[auth.BookshelfAccessToken].Value
			code := bookshelfCookies[auth.BookshelfTokenCode].Value
			parsedToken, err := auth.ParseJWT(log, accessToken)
			if err != nil {
				log.Errorf("could not parse access token: %v", err)
				apierr.APIErrorResponse(w, apierr.NewJWTTokenError(err.Error()))
				return
			}
			if ok, err := parsedToken.IsValid(); err != nil || !ok || !parsedToken.HasCorrectClaims(code) {
				log.Errorf("token not valid: valid - %+v error - %+v check - %t", parsedToken.Valid(), err, auth.CheckHash(parsedToken.Code, code))
				apierr.APIErrorResponse(w, apierr.NewJWTTokenError("invalid token"))
				return
			}
			ctx := request.AddAPIKeyToContext(r.Context(), parsedToken.Subject)
			req := r.WithContext(ctx)
			log.Info(parsedToken.Subject)
			next.ServeHTTP(w, req)
		})
	}
}

// AuthorizedSearch reads the JWT from the incoming request and redirects if the user is not authorized.
func AuthorizedSearch(log logs.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := os.Getenv("ALLOWED_URL_BASE") + "/webcli/error"
			cookies := r.Cookies()
			if len(cookies) < 1 {
				log.Error("no cookies in request")
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				return
			}
			bookshelfCookies, err := request.FindCookies(cookies, auth.BookshelfTokenCode, auth.BookshelfAccessToken)
			if err != nil {
				log.Errorf("could not find bookshelf cookies: %v", err)
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				return
			}
			accessToken := bookshelfCookies[auth.BookshelfAccessToken].Value
			code := bookshelfCookies[auth.BookshelfTokenCode].Value
			parsedToken, err := auth.ParseJWT(log, accessToken)
			if err != nil {
				log.Errorf("could not parse access token: %+v", err)
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				return
			}
			if !parsedToken.HasCorrectClaims(code) || len(accessToken) == 0 || len(code) == 0 {
				log.Errorf("token not valid: error - %+v check - %t", parsedToken, auth.CheckHash(parsedToken.Code, code))
				apierr.APIErrorResponse(w, apierr.NewJWTTokenError("invalid token"))
				return
			}
			refreshCode := ""
			ok, err := parsedToken.IsValid()
			if err != nil {
				log.Info("token not valid: %+v", err)
			}
			if !ok {
				log.Info("token refresh required on search")
				refreshCode = code
			}
			ctx := request.AddSearchKeysToContext(r.Context(), parsedToken.Subject, refreshCode)
			req := r.WithContext(ctx)
			log.Info(parsedToken.Subject)
			next.ServeHTTP(w, req)
		})
	}
}
