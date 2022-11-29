package auth

import (
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
)

var googleOAuth2Config = oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	Endpoint:     endpoints.Google,
	RedirectURL:  "http://localhost:8080/api/oauth/redirect",
	Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
}

var oidcConfig = &oidc.Config{
	ClientID: clientID,
}

type OAuth2Request struct {
	State, Nonce, AuthURL string
}

type GoogleIDTokenClaims struct {
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	PictureURL    string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Locale        string `json:"locale"`
}

type GoogleOIDCTokens struct {
	OAuth2Token   oauth2.Token
	IDTokenClaims GoogleIDTokenClaims
}
