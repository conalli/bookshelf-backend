package jwtauth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/conalli/bookshelf-backend/utils/apiErrors"
	"github.com/golang-jwt/jwt/v4"
)

var (
	key           = []byte(os.Getenv("SIGNING_SECRET"))
	headerTokenRe = regexp.MustCompile(`^Bearer\s([a-zA-Z0-9\.\-_]+)$`)
)

// CustomClaims represents the claims made in the JWT.
type CustomClaims struct {
	Name string
	jwt.StandardClaims
}

// NewToken creates a new token based on the CustomClaims and returns the token
// as a string signed with the secret.
func NewToken(name string) (string, apiErrors.ApiErr) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			NotBefore: time.Now().Add(5 * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "http://bookshelf-backend.jp",
		},
	})

	tkn, err := token.SignedString(key)
	if err != nil {
		log.Printf("error when trying to sign token %+v", token)
		return tkn, apiErrors.NewInternalServerError()
	}
	return tkn, nil
}

// Authorized reads the JWT from the incoming request and returns whether the user is authorized or not.
func Authorized(key string) func(w http.ResponseWriter, r *http.Request) bool {
	return func(w http.ResponseWriter, r *http.Request) bool {
		matches := headerTokenRe.FindStringSubmatch(r.Header.Get("Authorization"))
		if len(matches) < 2 {
			log.Println("error: failed to match authorization header")
			return false
		}
		token, err := jwt.ParseWithClaims(matches[1], &CustomClaims{}, func(t *jwt.Token) (interface{}, error) { return key, nil })
		if err != nil {
			log.Println("error: failed to parse with claims")
			return false
		}
		tkn, ok := token.Claims.(*CustomClaims)
		if !ok {
			log.Println("error: failed to convert token to CustomClaims")
			return false
		}
		if err = tkn.Valid(); err != nil {
			log.Println("error: token not valid")
			return false
		}
		fmt.Printf("%+v", tkn)
		return true
	}
}
