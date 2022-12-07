package auth_test

import (
	"os"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/golang-jwt/jwt/v4"
)

var signingKey = []byte(os.Getenv("SIGNING_SECRET"))

func TestNewToken(t *testing.T) {
	t.Parallel()
	tn := []string{
		"Name 1",
		"This is a test name",
		"abd",
		"12345",
		"Name",
	}

	for _, n := range tn {
		t.Run(n, func(t *testing.T) {
			tkn, err := auth.NewTokens(n, testutils.NewLogger())
			if err != nil {
				t.Fatalf("couldn't make a new token with name: %s", n)
			}
			token, e := jwt.ParseWithClaims(tkn["access_token"], &auth.CustomClaims{}, func(t *jwt.Token) (interface{}, error) { return signingKey, nil })
			if e != nil {
				t.Fatalf("couldn't parse token: %+v", e)
			}
			claimToken, ok := token.Claims.(*auth.CustomClaims)
			if !ok {
				t.Fatal("invalid custom claims")
			}
			if claimToken.Name != n {
				t.Fatalf("token.Name: %s not equal to name: %s", claimToken.Name, n)
			}
		})
	}
}
