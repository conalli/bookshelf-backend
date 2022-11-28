package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-playground/validator/v10"
	"golang.org/x/oauth2"
)

type Service interface {
	OAuthRequest(ctx context.Context) (OAuth2Request, error)
	OAuthRedirect(ctx context.Context, code, state string, stateCookie, nonceCookie *http.Cookie) (string, error)
}

type service struct {
	l logs.Logger
	v *validator.Validate
	p *oidc.Provider
}

func NewService(l logs.Logger, v *validator.Validate, p *oidc.Provider) *service {
	return &service{l, v, p}
}

func (s *service) OAuthRequest(ctx context.Context) (OAuth2Request, error) {
	state := ""
	nonce := ""
	url := googleOAuth2Config.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.AccessTypeOffline)
	return OAuth2Request{State: state, Nonce: nonce, AuthURL: url}, nil
}

func (s *service) OAuthRedirect(ctx context.Context, code, state string, stateCookie, nonceCookie *http.Cookie) (string, error) {
	if state != stateCookie.Value {
		s.l.Error("INCORRECT STATE FROM REDIRECT")
		return "", errors.NewBadRequestError("OMG")
	}
	oauth2Token, err := googleOAuth2Config.Exchange(ctx, code)
	if err != nil {
		s.l.Error(err)
		return "", err
	}
	verifier := s.p.Verifier(oidcConfig)
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		s.l.Error("NO id_token in token")
		return "", errors.NewInternalServerError()
	}
	IDToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		s.l.Error("Could not verify id_token:", err)
		return "", errors.NewInternalServerError()
	}
	if IDToken.Nonce != nonceCookie.Value {
		s.l.Error("Nonce did not match", IDToken.Nonce, nonceCookie.Value)
		return "", errors.NewBadRequestError("nonce did not match")
	}
	tokens := struct {
		OAuth2Token   oauth2.Token
		IDTokenClaims struct {
			Email, EmailVerified, FamilyName, GivenName, Locale, Name, Picture string
		}
	}{OAuth2Token: *oauth2Token}

	if err = IDToken.Claims(&tokens.IDTokenClaims); err != nil {
		s.l.Error("Could not parse id_token claims", err)
		return "", errors.NewInternalServerError()
	}
	j, err := json.Marshal(tokens)
	if err != nil {
		s.l.Error("Error marshaling JSON: ", err)
		return "", errors.NewInternalServerError()
	}

	s.l.Infof("IDToken: %+v", tokens)
	s.l.Infof("JSON: %+v", j)

	return "", nil
}
