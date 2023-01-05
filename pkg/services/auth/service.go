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
	NewUser(context.Context, request.SignUp) (string, error)
	NewOAuthUser(context.Context, GoogleIDTokenClaims) (string, error)
	NewRefreshToken(ctx context.Context, APIKey, refreshToken string) error
	GetRefreshTokenByAPIKey(ctx context.Context, APIKey string) (string, error)
}

type Service interface {
	SignUp(context.Context, request.SignUp) (*bookshelfTokens, errors.APIErr)
	LogIn(context.Context, request.LogIn) (*bookshelfTokens, errors.APIErr)
	OAuthRequest(ctx context.Context, authProvider, authType string) (OIDCRequest, errors.APIErr)
	OAuthRedirect(ctx context.Context, authProvider, authType, code, state string, cookies []*http.Cookie) (*bookshelfTokens, errors.APIErr)
	RefreshTokens(ctx context.Context, accessToken, code string) (*bookshelfTokens, errors.APIErr)
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
func (s *service) SignUp(ctx context.Context, requestData request.SignUp) (*bookshelfTokens, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Struct(requestData)
	if validateErr != nil {
		s.log.Errorf("Could not validate SIGN UP request: %v", validateErr)
		return nil, errors.NewBadRequestError("request format incorrect.")
	}
	userExists, err := s.db.UserAlreadyExists(ctx, requestData.Email)
	if err != nil {
		s.log.Errorf("error attempting to check if user exists: %v", err)
		return nil, errors.NewInternalServerError()
	}
	if userExists {
		s.log.Errorf("error creating new user; user with email %s already exists", requestData.Email)
		return nil, errors.NewBadRequestError("user already exists")
	}
	APIKey, err := s.db.NewUser(reqCtx, requestData)
	if err != nil {
		s.log.Errorf("couldnt create new user: %v", err)
		return nil, errors.NewInternalServerError()
	}
	tokens, err := NewTokens(s.log, APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return nil, errors.NewInternalServerError()
	}
	err = s.db.NewRefreshToken(ctx, APIKey, tokens.refreshToken)
	if err != nil {
		s.log.Error("could not save refresh token to db")
		return nil, errors.NewInternalServerError()
	}
	return tokens, nil
}

// Login takes in request data, checks the db and returns the username and apikey is successful.
func (s *service) LogIn(ctx context.Context, requestData request.LogIn) (*bookshelfTokens, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Struct(requestData)
	if validateErr != nil {
		s.log.Errorf("Could not validate LOG IN request: %v", validateErr)
		return nil, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.GetUserByEmail(reqCtx, requestData.Email)
	if err != nil || !CheckHashedPassword(user.Password, requestData.Password) {
		s.log.Errorf("could not login: %+v", err)
		return nil, errors.NewAPIError(http.StatusUnauthorized, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	tokens, err := NewTokens(s.log, user.APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return nil, errors.NewInternalServerError()
	}
	err = s.db.NewRefreshToken(ctx, user.APIKey, tokens.refreshToken)
	if err != nil {
		s.log.Error("could not save refresh token to db")
		return nil, errors.NewInternalServerError()
	}
	return tokens, nil
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

func (s *service) OAuthRedirect(ctx context.Context, authProvider, authType, code, state string, cookies []*http.Cookie) (*bookshelfTokens, errors.APIErr) {
	stateCookie := request.FilterCookies(cookies, "state")
	if stateCookie == nil {
		s.log.Error("no cookies in request")
		return nil, errors.NewBadRequestError("no cookies in auth request")
	}
	if state != stateCookie.Value {
		s.log.Error("state values did not match: %s - %s")
		return nil, errors.NewBadRequestError("invalid token")
	}
	switch authProvider {
	case "google":
		APIKey, apierr := s.GoogleOAuthRedirect(ctx, authType, code, cookies)
		if apierr != nil {
			return nil, apierr
		}
		tokens, err := NewTokens(s.log, APIKey)
		if err != nil {
			s.log.Error("could not create new tokens")
			return nil, errors.NewInternalServerError()
		}
		err = s.db.NewRefreshToken(ctx, APIKey, tokens.refreshToken)
		if err != nil {
			s.log.Error("could not save refresh token to db")
			return nil, errors.NewInternalServerError()
		}
		return tokens, nil
	default:
		s.log.Error("invalid auth provider in request url")
		return nil, errors.NewBadRequestError("invalid auth provider in request url")
	}
}

func (s *service) GoogleOAuthRedirect(ctx context.Context, authType, code string, cookies []*http.Cookie) (string, errors.APIErr) {
	nonceCookie := request.FilterCookies(cookies, "nonce")
	oauth2Token, err := newGoogleOAuth2Config(authType).Exchange(ctx, code)
	if err != nil {
		s.log.Error("could not exchange authorization code for token:", err)
		return "", errors.NewInternalServerError()
	}
	verifier := s.p.Verifier(googleOIDCConfig)
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		s.log.Error("no id_token in token")
		return "", errors.NewInternalServerError()
	}
	IDToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		s.log.Error("could not verify id_token:", err)
		return "", errors.NewInternalServerError()
	}
	if IDToken.Nonce != nonceCookie.Value {
		s.log.Errorf("nonces did not match: %s - %s", IDToken.Nonce, nonceCookie.Value)
		return "", errors.NewBadRequestError("invalid token")
	}
	oidcTokens := &GoogleOIDCTokens{OAuth2Token: *oauth2Token}
	if err = IDToken.Claims(&oidcTokens.IDTokenClaims); err != nil {
		s.log.Error("could not parse id_token claims", err)
		return "", errors.NewInternalServerError()
	}
	var APIKey string
	if authType == "login" {
		s.log.Info("login request")
		userInfo, err := s.db.GetUserByEmail(ctx, oidcTokens.IDTokenClaims.Email)
		if err != nil {
			s.log.Error(err)
			return "", errors.NewBadRequestError("couldnt find user with given email")
		}
		APIKey = userInfo.APIKey
	} else {
		s.log.Info("signup request")
		userExists, err := s.db.UserAlreadyExists(ctx, oidcTokens.IDTokenClaims.Email)
		if err != nil {
			s.log.Errorf("error attempting to check if user exists: %v", err)
			return "", errors.NewInternalServerError()
		}
		if userExists {
			s.log.Errorf("error creating new user; user with email %s already exists", oidcTokens.IDTokenClaims.Email)
			return "", errors.NewBadRequestError("user already exists")
		}
		userInfo, err := s.db.NewOAuthUser(ctx, oidcTokens.IDTokenClaims)
		if err != nil {
			s.log.Errorf("couldnt create user from id token: %v", err)
			return "", errors.NewInternalServerError()
		}
		APIKey = userInfo
	}
	return APIKey, nil
}

func (s *service) RefreshTokens(ctx context.Context, accessToken, code string) (*bookshelfTokens, errors.APIErr) {
	tkn, err := ParseJWT(s.log, accessToken, code)
	if err != nil {
		s.log.Error("could not parse jwt from cookie")
		return nil, errors.NewJWTTokenError("could not parse token")
	}
	APIKey := tkn.Subject
	token, err := s.db.GetRefreshTokenByAPIKey(ctx, APIKey)
	if err != nil {
		s.log.Error("could not get refresh token from db")
		if err == errors.ErrInternalServerError {
			return nil, errors.NewInternalServerError()
		}
		return nil, errors.NewAPIError(http.StatusNotFound, err.Error(), "no refresh token")
	}
	_, err = ParseJWT(s.log, token, code)
	if err != nil {
		s.log.Error("parsed refresh token invalid")
		return nil, errors.NewBadRequestError("invalid refresh token")
	}
	tokens, err := NewTokens(s.log, APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return nil, errors.NewInternalServerError()
	}
	err = s.db.NewRefreshToken(ctx, APIKey, tokens.refreshToken)
	if err != nil {
		s.log.Error("could not save refresh token to db")
		return nil, errors.NewInternalServerError()
	}
	return tokens, nil
}
