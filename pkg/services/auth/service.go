package auth

import (
	"context"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
)

type Service interface {
	OAuthRequest(ctx context.Context, authType string) (OIDCRequest, error)
	OAuthRedirect(ctx context.Context, code, state string, stateCookie, nonceCookie *http.Cookie) (*GoogleOIDCTokens, errors.APIErr)
}

type service struct {
	l logs.Logger
	v *validator.Validate
	p *oidc.Provider
}

func NewService(l logs.Logger, v *validator.Validate, p *oidc.Provider) *service {
	return &service{l, v, p}
}

func (s *service) OAuthRequest(ctx context.Context, authType string) (OIDCRequest, error) {
	state := authType + ""
	nonce := ""
	url := googleOAuth2Config.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.AccessTypeOffline)
	return OIDCRequest{State: state, Nonce: nonce, AuthURL: url}, nil
}

func (s *service) OAuthRedirect(ctx context.Context, code, state string, stateCookie, nonceCookie *http.Cookie) (*GoogleOIDCTokens, errors.APIErr) {
	if state != stateCookie.Value {
		s.l.Error("state values did not match: %s - %s")
		return nil, errors.NewBadRequestError("invalid token")
	}
	oauth2Token, err := googleOAuth2Config.Exchange(ctx, code)
	if err != nil {
		s.l.Error("could not exchange authorization code for token:", err)
		return nil, errors.NewInternalServerError()
	}
	verifier := s.p.Verifier(googleOIDCConfig)
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		s.l.Error("no id_token in token")
		return nil, errors.NewInternalServerError()
	}
	IDToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		s.l.Error("could not verify id_token:", err)
		return nil, errors.NewInternalServerError()
	}
	if IDToken.Nonce != nonceCookie.Value {
		s.l.Errorf("nonces did not match: %s - %s", IDToken.Nonce, nonceCookie.Value)
		return nil, errors.NewBadRequestError("invalid token")
	}
	tokens := &GoogleOIDCTokens{OAuth2Token: *oauth2Token}
	if err = IDToken.Claims(&tokens.IDTokenClaims); err != nil {
		s.l.Error("could not parse id_token claims", err)
		return nil, errors.NewInternalServerError()
	}
	s.l.Infof("Tokens: %+v", tokens)
	return tokens, nil
}
