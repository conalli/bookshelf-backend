package jwtauth

import (
	"net/http"
	"os"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var signingKey = []byte(os.Getenv("SIGNING_SECRET"))

// CustomClaims represents the claims made in the JWT.
type CustomClaims struct {
	Name string
	jwt.RegisteredClaims
}

// NewTokens creates a new token based on the CustomClaims and returns the token
// as a string signed with the secret.
func NewTokens(name string, log logs.Logger) (map[string]string, errors.APIErr) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "http://bookshelf-backend.jp",
			Subject:   name,
		},
	})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "http://bookshelf-backend.jp",
		Subject:   name,
	})
	tkn, tknErr := token.SignedString(signingKey)
	ref, refErr := refresh.SignedString(signingKey)
	if tknErr != nil || refErr != nil {
		log.Errorf("error when trying to sign tokens %+v", token)
		return nil, errors.NewInternalServerError()
	}
	tkns := map[string]string{
		"access_token":  tkn,
		"refresh_token": ref,
	}
	return tkns, nil
}

// Authorized reads the JWT from the incoming request and returns whether the user is authorized or not.
func Authorized(next http.HandlerFunc, log logs.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["APIKey"]
		cookies := r.Cookies()
		if len(cookies) < 1 {
			log.Error("no cookies in request")
			errors.APIErrorResponse(w, errors.NewBadRequestError("no cookies in request"))
			return
		}
		bookshelfCookie := FilterCookies("bookshelfjwt", cookies)
		refreshCookie := FilterCookies("bookshelfrefresh", cookies)
		if bookshelfCookie == nil && refreshCookie == nil {
			log.Error("did not find bookshelf cookies")
			errors.APIErrorResponse(w, errors.NewBadRequestError("did not find bookshelf cookies"))
			return
		}
		if bookshelfCookie != nil {
			token, err := jwt.ParseWithClaims(bookshelfCookie.Value, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) { return signingKey, nil })
			tkn, ok := token.Claims.(*CustomClaims)
			if !ok || err != nil {
				log.Error("failed to convert token to CustomClaims")
				errors.APIErrorResponse(w, errors.NewJWTClaimsError("failed to convert token to CustomClaims"))
				return
			}
			if err = tkn.Valid(); err != nil || len(name) == 0 || tkn.Name != name || tkn.Subject != name {
				log.Error("token not valid")
				errors.APIErrorResponse(w, errors.NewJWTTokenError("error: token not valid"))
				return
			}
		} else {
			refresh, err := jwt.ParseWithClaims(refreshCookie.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return signingKey, nil })
			if err != nil {
				log.Errorf("failed to parse refresh token with claims: %+v\n", err)
				errors.APIErrorResponse(w, errors.NewJWTClaimsError("failed to parse with claims"))
				return
			}
			ref, ok := refresh.Claims.(*jwt.RegisteredClaims)
			if !ok {
				log.Error("failed to convert refresh token")
				errors.APIErrorResponse(w, errors.NewJWTClaimsError("failed to convert token"))
				return
			}
			if err = ref.Valid(); err != nil || len(name) == 0 || ref.Subject != name {
				log.Error("refresh token not valid")
				errors.APIErrorResponse(w, errors.NewJWTTokenError("token not valid"))
				return
			}
			tokens, err := NewTokens(name, log)
			if err != nil {
				log.Errorf("could not create new bookshelf cookies: %v", err)
				errors.APIErrorResponse(w, errors.NewInternalServerError())
				return
			}
			accessToken := http.Cookie{
				Name:     "bookshelfjwt",
				Value:    tokens["access_token"],
				Expires:  time.Now().Add(15 * time.Minute),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			}
			refreshToken := http.Cookie{
				Name:     "bookshelfrefresh",
				Value:    tokens["refresh_token"],
				Expires:  time.Now().Add(24 * time.Hour),
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			}
			log.Info("successfully returned tokens as cookies")
			http.SetCookie(w, &accessToken)
			http.SetCookie(w, &refreshToken)
		}
		next(w, r)
	}
}

// FilterCookies looks through all cookies and returns one with given name.
func FilterCookies(name string, cookies []*http.Cookie) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}
