package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
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
		return nil, apierr.ErrInternalServerError
	}
	code := jwtid.String()
	codeHash, err := HashPassword(code)
	if err != nil {
		log.Error("could not hash jwt code")
		return nil, apierr.ErrInternalServerError
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTCustomClaims{
		Code: codeHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Minute)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("SERVER_URL_BASE"),
			Subject:   APIKey,
		},
	})
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTCustomClaims{
		Code: codeHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("SERVER_URL_BASE"),
			Subject:   APIKey,
		},
	})
	access, tknErr := token.SignedString(signingKey)
	ref, refErr := refresh.SignedString(signingKey)
	if tknErr != nil || refErr != nil {
		log.Errorf("error when trying to sign tokens %+v", token)
		return nil, apierr.ErrInternalServerError
	}
	tokens := &bookshelfTokens{code, access, ref}
	return tokens, nil
}

func (t *bookshelfTokens) NewTokenCookies(log logs.Logger) []*http.Cookie {
	now := time.Now()
	codeExpires := now.Add(24 * time.Hour)
	accessExpires := now.Add(20 * time.Minute)
	path := "/"
	secure := true
	httpOnly := true
	sameSite := http.SameSiteNoneMode

	codeCookie := &http.Cookie{
		Name:     BookshelfTokenCode,
		Value:    t.code,
		Path:     path,
		Expires:  codeExpires,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
		MaxAge:   24 * 60 * 60,
	}

	accessCookie := &http.Cookie{
		Name:     BookshelfAccessToken,
		Value:    t.accessToken,
		Path:     path,
		Expires:  accessExpires,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
		MaxAge:   20 * 60,
	}

	return []*http.Cookie{codeCookie, accessCookie}
}

func AddCookiesToResponse(w http.ResponseWriter, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
}

func RemoveBookshelfCookies(w http.ResponseWriter) {
	path := "/"
	expires := time.Now().Add(-100 * time.Hour)
	secure := true
	httpOnly := true
	sameSite := http.SameSiteStrictMode
	maxAge := -1
	codeCookie := &http.Cookie{
		Name:     BookshelfTokenCode,
		Value:    "",
		Path:     path,
		Expires:  expires,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
		MaxAge:   maxAge,
	}
	accessCookie := &http.Cookie{
		Name:     BookshelfAccessToken,
		Value:    "",
		Path:     path,
		Expires:  expires,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
		MaxAge:   maxAge,
	}
	http.SetCookie(w, codeCookie)
	http.SetCookie(w, accessCookie)
}

func ParseJWT(log logs.Logger, token, code string) (*JWTCustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTCustomClaims{}, func(t *jwt.Token) (interface{}, error) { return signingKey, nil })
	tkn, ok := parsedToken.Claims.(*JWTCustomClaims)
	if !ok || err != nil {
		log.Error("failed to convert token to JWTCustomClaims")
		return nil, apierr.ErrInvalidJWTToken
	}
	if err = tkn.Valid(); err != nil || !CheckHashedPassword(tkn.Code, code) {
		log.Errorf("token not valid: error - %v check - %t", err, CheckHashedPassword(tkn.Code, code))
		return nil, apierr.ErrInvalidJWTClaims
	}
	return tkn, nil
}
