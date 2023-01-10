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
	NewUser(context.Context, accounts.User) (string, error)
	NewRefreshToken(ctx context.Context, APIKey, refreshToken string) error
	GetRefreshTokenByAPIKey(ctx context.Context, APIKey string) (string, error)
}

type Cache interface {
	AddUser(ctx context.Context, userKey string, user accounts.User) (int64, error)
	DeleteUser(ctx context.Context, userKey string) (int64, error)
}

type Service interface {
	SignUp(context.Context, request.SignUp) (AuthUser, errors.APIErr)
	LogIn(context.Context, request.LogIn) (AuthUser, errors.APIErr)
	OAuthRequest(ctx context.Context, authProvider, authType string) (OIDCRequest, errors.APIErr)
	OAuthRedirect(ctx context.Context, authProvider, authType, code, state string, cookies []*http.Cookie) (*bookshelfTokens, errors.APIErr)
	RefreshTokens(ctx context.Context, accessToken, code string) (*bookshelfTokens, errors.APIErr)
}

type service struct {
	log      logs.Logger
	validate *validator.Validate
	p        *oidc.Provider
	db       Repository
	cache    Cache
}

func NewService(l logs.Logger, v *validator.Validate, p *oidc.Provider, db Repository, c Cache) *service {
	return &service{l, v, p, db, c}
}

type AuthUser struct {
	User   accounts.User
	Tokens *bookshelfTokens
}

// SignUp returns the url of a given cmd.
func (s *service) SignUp(ctx context.Context, requestData request.SignUp) (AuthUser, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Struct(requestData)
	if err != nil {
		s.log.Errorf("Could not validate SIGN UP request: %v", err)
		return AuthUser{}, errors.NewBadRequestError("request format incorrect.")
	}
	userExists, err := s.db.UserAlreadyExists(ctx, requestData.Email)
	if err != nil {
		s.log.Errorf("error attempting to check if user exists: %v", err)
		return AuthUser{}, errors.NewInternalServerError()
	}
	if userExists {
		s.log.Errorf("error creating new user; user with email %s already exists", requestData.Email)
		return AuthUser{}, errors.NewBadRequestError("user already exists")
	}
	APIKey, err := GenerateAPIKey()
	if err != nil {
		s.log.Error("could not generate uuid")
		return AuthUser{}, errors.NewInternalServerError()
	}
	hashedPassword, err := HashPassword(requestData.Password)
	if err != nil {
		s.log.Error("could not hash password")
		return AuthUser{}, errors.NewInternalServerError()
	}
	user := accounts.User{
		Email:    requestData.Email,
		Password: hashedPassword,
		APIKey:   APIKey,
		Provider: ProviderBookshelf,
		Cmds:     map[string]string{},
		Teams:    map[string]string{},
	}
	userID, err := s.db.NewUser(reqCtx, user)
	if err != nil {
		s.log.Errorf("couldnt create new user: %v", err)
		return AuthUser{}, errors.NewInternalServerError()
	}
	user.ID = userID
	tokens, err := NewTokens(s.log, user.APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return AuthUser{}, errors.NewInternalServerError()
	}
	err = s.db.NewRefreshToken(ctx, user.APIKey, tokens.refreshToken)
	if err != nil {
		s.log.Error("could not save refresh token to db")
		return AuthUser{}, errors.NewInternalServerError()
	}
	_, err = s.cache.AddUser(ctx, user.APIKey, user)
	if err != nil {
		s.log.Error("could not add user to cache")
	}
	authUser := AuthUser{
		User:   user,
		Tokens: tokens,
	}
	return authUser, nil
}

// Login takes in request data, checks the db and returns the username and apikey is successful.
func (s *service) LogIn(ctx context.Context, requestData request.LogIn) (AuthUser, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Struct(requestData)
	if err != nil {
		s.log.Errorf("Could not validate LOG IN request: %v", err)
		return AuthUser{}, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.GetUserByEmail(reqCtx, requestData.Email)
	if err != nil || !CheckHashedPassword(user.Password, requestData.Password) {
		s.log.Errorf("could not login: %+v", err)
		return AuthUser{}, errors.NewAPIError(http.StatusUnauthorized, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	tokens, err := NewTokens(s.log, user.APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return AuthUser{}, errors.NewInternalServerError()
	}
	err = s.db.NewRefreshToken(ctx, user.APIKey, tokens.refreshToken)
	if err != nil {
		s.log.Error("could not save refresh token to db")
	}
	_, err = s.cache.AddUser(ctx, user.APIKey, user)
	if err != nil {
		s.log.Error("could not add user to cache")
	}
	authUser := AuthUser{
		User:   user,
		Tokens: tokens,
	}
	return authUser, nil
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
		s.log.Error("no state cookie in OAuth redirect")
		return nil, errors.NewBadRequestError("no cookies in auth request")
	}
	if state != stateCookie.Value {
		s.log.Error("state values did not match: %s - %s")
		return nil, errors.NewBadRequestError("invalid token")
	}
	var APIKey string
	switch authProvider {
	case ProviderGoogle:
		key, apierr := s.googleOIDCRedirect(ctx, authType, code, cookies)
		if apierr != nil {
			return nil, apierr
		}
		APIKey = key
	default:
		s.log.Error("invalid auth provider in request url")
		return nil, errors.NewBadRequestError("invalid auth provider in request url")
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
