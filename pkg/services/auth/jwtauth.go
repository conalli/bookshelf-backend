package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	BookshelfTokenCode    string = "bookshelf_token_code"
	BookshelfAccessToken  string = "bookshelf_access_token"
	BookshelfRefreshToken string = "bookshelf_refresh_token"
)

var signingKey = []byte(os.Getenv("SIGNING_SECRET"))

// CustomClaims represents the claims made in the JWT.
type JWTCustomClaims struct {
	Code string
	jwt.RegisteredClaims
}

type bookshelfTokens struct {
	code, accessToken, refreshToken string
}

func (b *bookshelfTokens) Code() string {
	return b.code
}

func (b *bookshelfTokens) AccessToken() string {
	return b.accessToken
}

// NewTokens creates a new token based on the CustomClaims and returns the token
// as a string signed with the secret.
func NewTokens(log logs.Logger, APIKey string) (*bookshelfTokens, error) {
	jwtid, err := uuid.NewRandom()
	if err != nil {
		log.Error("could not generate uuid for jwt")
		return nil, errors.ErrInternalServerError
	}
	code := jwtid.String()
	codeHash, err := HashPassword(code)
	if err != nil {
		log.Error("could not hash jwt code")
		return nil, errors.ErrInternalServerError
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTCustomClaims{
		Code: codeHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "https://localhost:8080/api",
			Subject:   APIKey,
		},
	})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTCustomClaims{
		Code: codeHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "https://localhost:8080/api",
			Subject:   APIKey,
		},
	})
	access, tknErr := token.SignedString(signingKey)
	ref, refErr := refresh.SignedString(signingKey)
	if tknErr != nil || refErr != nil {
		log.Errorf("error when trying to sign tokens %+v", token)
		return nil, errors.ErrInternalServerError
	}
	tokens := &bookshelfTokens{code, access, ref}
	return tokens, nil
}

func (t *bookshelfTokens) NewTokenCookies(log logs.Logger) []*http.Cookie {
	now := time.Now()
	codeExpires := now.Add(24 * time.Hour)
	accessExpires := now.Add(15 * time.Minute)
	path := "/"
	secure := true
	httpOnly := false
	sameSite := http.SameSiteNoneMode

	codeCookie := &http.Cookie{
		Name:     BookshelfTokenCode,
		Value:    t.code,
		Path:     path,
		Expires:  codeExpires,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	}

	accessCookie := &http.Cookie{
		Name:     BookshelfAccessToken,
		Value:    t.accessToken,
		Path:     path,
		Expires:  accessExpires,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	}

	return []*http.Cookie{codeCookie, accessCookie}
}

func AddCookiesToResponse(w http.ResponseWriter, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
}

func ParseJWT(log logs.Logger, token, code string) (*JWTCustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTCustomClaims{}, func(t *jwt.Token) (interface{}, error) { return signingKey, nil })
	tkn, ok := parsedToken.Claims.(*JWTCustomClaims)
	if !ok || err != nil {
		log.Error("failed to convert token to JWTCustomClaims")
		return nil, errors.ErrInvalidJWTToken
	}
	if err = tkn.Valid(); err != nil || !CheckHashedPassword(tkn.Code, code) {
		log.Errorf("token not valid: error - %v check - %t", err, CheckHashedPassword(tkn.Code, code))
		return nil, errors.ErrInvalidJWTClaims
	}
	return tkn, nil
}
