package auth

import (
	"context"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Repository interface {
	UserAlreadyExists(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (accounts.User, error)
	NewUser(context.Context, request.SignUp) (accounts.User, error)
	NewOAuthUser(context.Context, GoogleIDTokenClaims) (accounts.User, error)
}

type Service interface {
	SignUp(context.Context, request.SignUp) (accounts.User, errors.APIErr)
	LogIn(context.Context, request.LogIn) (accounts.User, errors.APIErr)
	OAuthRequest(ctx context.Context, authProvider, authType string) (OIDCRequest, errors.APIErr)
	OAuthRedirect(ctx context.Context, authProvider, authType, code, state string, cookies []*http.Cookie) (accounts.User, errors.APIErr)
}

type service struct {
	log      logs.Logger
	validate *validator.Validate
	p        *oidc.Provider
	db       Repository
}

func NewService(l logs.Logger, v *validator.Validate, p *oidc.Provider, db Repository) *service {
	return &service{l, v, p, db}
}

// SignUp returns the url of a given cmd.
func (s *service) SignUp(ctx context.Context, requestData request.SignUp) (accounts.User, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Struct(requestData)
	if validateErr != nil {
		s.log.Errorf("Could not validate SIGN UP request: %v", validateErr)
		return accounts.User{}, errors.NewBadRequestError("request format incorrect.")
	}
	userExists, err := s.db.UserAlreadyExists(ctx, requestData.Email)
	if err != nil {
		s.log.Errorf("error attempting to check if user exists: %v", err)
		return accounts.User{}, errors.NewInternalServerError()
	}
	if userExists {
		s.log.Errorf("error creating new user; user with email %s already exists", requestData.Email)
		return accounts.User{}, errors.NewBadRequestError("user already exists")
	}
	user, err := s.db.NewUser(reqCtx, requestData)
	if err != nil {
		s.log.Errorf("couldnt create new user: %v", err)
		return accounts.User{}, errors.NewInternalServerError()
	}
	return user, nil
}

// Login takes in request data, checks the db and returns the username and apikey is successful.
func (s *service) LogIn(ctx context.Context, requestData request.LogIn) (accounts.User, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Struct(requestData)
	if validateErr != nil {
		s.log.Errorf("Could not validate LOG IN request: %v", validateErr)
		return accounts.User{}, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.GetUserByEmail(reqCtx, requestData.Email)
	if err != nil || !CheckHashedPassword(user.Password, requestData.Password) {
		s.log.Errorf("could not login: %+v", err)
		return accounts.User{}, errors.NewAPIError(http.StatusUnauthorized, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	return user, nil
}

func (s *service) OAuthRequest(ctx context.Context, authProvider, authType string) (OIDCRequest, errors.APIErr) {
	stateID, err := uuid.NewRandom()
	if err != nil {
		return OIDCRequest{}, errors.NewInternalServerError()
	}
	nonceID, err := uuid.NewRandom()
	if err != nil {
		return OIDCRequest{}, errors.NewInternalServerError()
	}
	state := authType + stateID.String()
	nonce := nonceID.String()
	url := newGoogleOAuth2Config(authType).AuthCodeURL(state, oidc.Nonce(nonce), oauth2.AccessTypeOffline)
	s.log.Infof("OAuth redirect url: %s", url)
	return OIDCRequest{State: state, Nonce: nonce, AuthURL: url}, nil
}

func (s *service) OAuthRedirect(ctx context.Context, authProvider, authType, code, state string, cookies []*http.Cookie) (accounts.User, errors.APIErr) {
	stateCookie := request.FilterCookies("state", cookies)
	if stateCookie == nil {
		s.log.Error("no cookies in request")
		return accounts.User{}, errors.NewBadRequestError("no cookies in auth request")
	}
	if state != stateCookie.Value {
		s.log.Error("state values did not match: %s - %s")
		return accounts.User{}, errors.NewBadRequestError("invalid token")
	}
	switch authProvider {
	case "google":
		user, err := s.GoogleOAuthRedirect(ctx, authType, code, cookies)
		if err != nil {
			return accounts.User{}, err
		}
		return user, nil
	default:
		s.log.Error("invalid auth provider in request url")
		return accounts.User{}, errors.NewBadRequestError("invalid auth provider in request url")
	}
}

func (s *service) GoogleOAuthRedirect(ctx context.Context, authType, code string, cookies []*http.Cookie) (accounts.User, errors.APIErr) {
	nonceCookie := request.FilterCookies("nonce", cookies)
	oauth2Token, err := newGoogleOAuth2Config(authType).Exchange(ctx, code)
	if err != nil {
		s.log.Error("could not exchange authorization code for token:", err)
		return accounts.User{}, errors.NewInternalServerError()
	}
	verifier := s.p.Verifier(googleOIDCConfig)
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		s.log.Error("no id_token in token")
		return accounts.User{}, errors.NewInternalServerError()
	}
	IDToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		s.log.Error("could not verify id_token:", err)
		return accounts.User{}, errors.NewInternalServerError()
	}
	if IDToken.Nonce != nonceCookie.Value {
		s.log.Errorf("nonces did not match: %s - %s", IDToken.Nonce, nonceCookie.Value)
		return accounts.User{}, errors.NewBadRequestError("invalid token")
	}
	tokens := &GoogleOIDCTokens{OAuth2Token: *oauth2Token}
	if err = IDToken.Claims(&tokens.IDTokenClaims); err != nil {
		s.log.Error("could not parse id_token claims", err)
		return accounts.User{}, errors.NewInternalServerError()
	}
	s.log.Infof("Tokens: %+v", tokens)
	var user accounts.User
	if authType == "login" {
		s.log.Info("login request")
		userInfo, err := s.db.GetUserByEmail(ctx, tokens.IDTokenClaims.Email)
		if err != nil {
			s.log.Error(err)
			return accounts.User{}, errors.NewBadRequestError("couldnt find user with given email")
		}
		user = userInfo
	} else {
		s.log.Info("signup request")
		userExists, err := s.db.UserAlreadyExists(ctx, tokens.IDTokenClaims.Email)
		if err != nil {
			s.log.Errorf("error attempting to check if user exists: %v", err)
			return accounts.User{}, errors.NewInternalServerError()
		}
		if userExists {
			s.log.Errorf("error creating new user; user with email %s already exists", tokens.IDTokenClaims.Email)
			return accounts.User{}, errors.NewBadRequestError("user already exists")
		}
		userInfo, err := s.db.NewOAuthUser(ctx, tokens.IDTokenClaims)
		if err != nil {
			s.log.Errorf("couldnt create user from id token: %v", err)
			return accounts.User{}, errors.NewInternalServerError()
		}
		user = userInfo
	}
	return user, nil
}
