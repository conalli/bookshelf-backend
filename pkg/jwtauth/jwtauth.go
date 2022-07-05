package jwtauth

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

var signingKey = []byte(os.Getenv("SIGNING_SECRET"))

// CustomClaims represents the claims made in the JWT.
type CustomClaims struct {
	Name string
	jwt.RegisteredClaims
}

// NewToken creates a new token based on the CustomClaims and returns the token
// as a string signed with the secret.
func NewToken(name string) (string, errors.ApiErr) {
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

	tkn, err := token.SignedString(signingKey)
	if err != nil {
		log.Printf("error when trying to sign token %+v", token)
		return tkn, errors.NewInternalServerError()
	}
	return tkn, nil
}

// Authorized reads the JWT from the incoming request and returns whether the user is authorized or not.
// TODO: improve validation
func Authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["APIKey"]
		cookies := r.Cookies()
		if len(cookies) < 1 {
			log.Println("error: no cookies in request")
			errors.APIErrorResponse(w, errors.NewBadRequestError("error: no cookies in request"))
			return
		}
		bookshelfCookie := FilterCookies("bookshelfjwt", cookies)
		if bookshelfCookie == nil {
			log.Println("error: did not find bookshelf cookie")
			errors.APIErrorResponse(w, errors.NewBadRequestError("error: did not find bookshelf cookie"))
			return
		}
		token, err := jwt.ParseWithClaims(bookshelfCookie.Value, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) { return signingKey, nil })
		if err != nil {
			log.Println("error: failed to parse with claims")
			log.Printf("error: %+v\n", err)
			errors.APIErrorResponse(w, errors.NewJWTClaimsError("error: failed to parse with claims"))
			return
		}
		tkn, ok := token.Claims.(*CustomClaims)
		if !ok {
			log.Println("error: failed to convert token to CustomClaims")
			errors.APIErrorResponse(w, errors.NewJWTClaimsError("error: failed to convert token to CustomClaims"))
			return
		}
		if err = tkn.Valid(); err != nil || tkn.Name != name || tkn.Subject != name {
			log.Println("error: token not valid")
			errors.APIErrorResponse(w, errors.NewJWTTokenError("error: token not valid"))
			return
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
