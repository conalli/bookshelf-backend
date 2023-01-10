package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
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

func (s *service) googleOIDCRedirect(ctx context.Context, authType, code string, cookies []*http.Cookie) (string, apierr.Error) {
	nonceCookie := request.FilterCookies(cookies, "nonce")
	oauth2Token, err := newGoogleOAuth2Config(authType).Exchange(ctx, code)
	if err != nil {
		s.log.Error("could not exchange authorization code for token:", err)
		return "", apierr.NewInternalServerError()
	}
	verifier := s.p.Verifier(googleOIDCConfig)
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		s.log.Error("no id_token in token")
		return "", apierr.NewInternalServerError()
	}
	IDToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		s.log.Error("could not verify id_token:", err)
		return "", apierr.NewInternalServerError()
	}
	if IDToken.Nonce != nonceCookie.Value {
		s.log.Errorf("nonces did not match: %s - %s", IDToken.Nonce, nonceCookie.Value)
		return "", apierr.NewBadRequestError("invalid token")
	}
	oidcTokens := &GoogleOIDCTokens{OAuth2Token: *oauth2Token}
	if err = IDToken.Claims(&oidcTokens.IDTokenClaims); err != nil {
		s.log.Error("could not parse id_token claims", err)
		return "", apierr.NewInternalServerError()
	}
	var APIKey string
	switch authType {
	case AuthTypeLogIn:
		key, apierr := s.googleOIDCLogIn(ctx, oidcTokens.IDTokenClaims.Email)
		if err != nil {
			return "", apierr
		}
		APIKey = key
	case AuthTypeSignUp:
		key, apierr := s.googleOIDCSignUp(ctx, oidcTokens.IDTokenClaims)
		if err != nil {
			return "", apierr
		}
		APIKey = key
	default:
		s.log.Error("invalid auth type in google redirect")
		return "", apierr.NewBadRequestError("invalid auth type in request")
	}
	return APIKey, nil
}

func (s *service) googleOIDCLogIn(ctx context.Context, email string) (string, apierr.Error) {
	s.log.Info("google login request")
	userInfo, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Error(err)
		return "", apierr.NewBadRequestError("couldnt find user with given email")
	}
	return userInfo.APIKey, nil
}

func (s *service) googleOIDCSignUp(ctx context.Context, claims GoogleIDTokenClaims) (string, apierr.Error) {
	s.log.Info("signup request")
	userExists, err := s.db.UserAlreadyExists(ctx, claims.Email)
	if err != nil {
		s.log.Errorf("error attempting to check if user exists: %v", err)
		return "", apierr.NewInternalServerError()
	}
	if userExists {
		s.log.Errorf("error creating new user; user with email %s already exists", claims.Email)
		return "", apierr.NewBadRequestError("user already exists")
	}
	newAPIKey, err := GenerateAPIKey()
	if err != nil {
		s.log.Error("could not generate uuid")
		return "", apierr.NewInternalServerError()
	}
	user := accounts.User{
		APIKey:        newAPIKey,
		Name:          claims.Name,
		GivenName:     claims.GivenName,
		FamilyName:    claims.FamilyName,
		PictureURL:    claims.PictureURL,
		Email:         claims.Email,
		EmailVerified: claims.EmailVerified,
		Locale:        claims.Locale,
		Provider:      ProviderGoogle,
		Cmds:          map[string]string{},
		Teams:         map[string]string{},
	}
	_, err = s.db.NewUser(ctx, user)
	if err != nil {
		s.log.Errorf("couldnt create user from id token: %v", err)
		return "", apierr.NewInternalServerError()
	}
	return user.APIKey, nil
}
