package auth

import (
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var (
	googleClientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
	googleOIDCConfig   = &oidc.Config{
		ClientID: googleClientID,
	}
)

type GoogleIDTokenClaims struct {
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	PictureURL    string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func newGoogleOAuth2Config(authType string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Endpoint:     endpoints.Google,
		RedirectURL:  fmt.Sprintf("http://localhost:8080/api/auth/redirect/google/%s", authType),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
}

type GoogleOIDCTokens struct {
	OAuth2Token   oauth2.Token
	IDTokenClaims GoogleIDTokenClaims
}
