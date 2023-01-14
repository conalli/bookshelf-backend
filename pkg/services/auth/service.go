package auth

import (
	"context"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
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
	DeleteRefreshToken(ctx context.Context, APIKey string) (int64, error)
}

type Cache interface {
	AddUser(ctx context.Context, userKey string, user accounts.User) (int64, error)
	DeleteUser(ctx context.Context, userKey string) (int64, error)
}

type Service interface {
	SignUp(context.Context, request.SignUp) (AuthUser, apierr.Error)
	LogIn(context.Context, request.LogIn) (AuthUser, apierr.Error)
	OAuthRequest(ctx context.Context, authProvider, authType string) (OIDCRequest, apierr.Error)
	OAuthRedirect(ctx context.Context, authProvider, authType, code, state string, cookies []*http.Cookie) (*BookshelfTokens, apierr.Error)
	RefreshTokens(ctx context.Context, accessToken, code string) (*BookshelfTokens, apierr.Error)
	LogOut(ctx context.Context, APIKey string) apierr.Error
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
	Tokens *BookshelfTokens
}

// SignUp returns the url of a given cmd.
func (s *service) SignUp(ctx context.Context, requestData request.SignUp) (AuthUser, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Struct(requestData)
	if err != nil {
		s.log.Errorf("Could not validate SIGN UP request: %v", err)
		return AuthUser{}, apierr.NewBadRequestError("request format incorrect.")
	}
	userExists, err := s.db.UserAlreadyExists(ctx, requestData.Email)
	if err != nil {
		s.log.Errorf("error attempting to check if user exists: %v", err)
		return AuthUser{}, apierr.NewInternalServerError()
	}
	if userExists {
		s.log.Errorf("error creating new user; user with email %s already exists", requestData.Email)
		return AuthUser{}, apierr.NewBadRequestError("user already exists")
	}
	APIKey, err := GenerateAPIKey()
	if err != nil {
		s.log.Error("could not generate uuid")
		return AuthUser{}, apierr.NewInternalServerError()
	}
	hashedPassword, err := Hash(requestData.Password)
	if err != nil {
		s.log.Error("could not hash password")
		return AuthUser{}, apierr.NewInternalServerError()
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
		return AuthUser{}, apierr.NewInternalServerError()
	}
	user.ID = userID
	tokens, err := NewTokens(s.log, user.APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return AuthUser{}, apierr.NewInternalServerError()
	}
	err = s.db.NewRefreshToken(ctx, user.APIKey, tokens.refreshToken)
	if err != nil {
		s.log.Error("could not save refresh token to db")
		return AuthUser{}, apierr.NewInternalServerError()
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
func (s *service) LogIn(ctx context.Context, requestData request.LogIn) (AuthUser, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Struct(requestData)
	if err != nil {
		s.log.Errorf("Could not validate LOG IN request: %v", err)
		return AuthUser{}, apierr.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.GetUserByEmail(reqCtx, requestData.Email)
	if err != nil || !CheckHash(user.Password, requestData.Password) {
		s.log.Errorf("could not login: %+v", err)
		return AuthUser{}, apierr.NewAPIError(http.StatusUnauthorized, apierr.ErrWrongCredentials, "error: name or password incorrect")
	}
	tokens, err := NewTokens(s.log, user.APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return AuthUser{}, apierr.NewInternalServerError()
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

func (s *service) OAuthRequest(ctx context.Context, authProvider, authType string) (OIDCRequest, apierr.Error) {
	stateID, err := uuid.NewRandom()
	if err != nil {
		return OIDCRequest{}, apierr.NewInternalServerError()
	}
	nonceID, err := uuid.NewRandom()
	if err != nil {
		return OIDCRequest{}, apierr.NewInternalServerError()
	}
	state := authType + stateID.String()
	nonce := nonceID.String()
	url := newGoogleOAuth2Config(authType).AuthCodeURL(state, oidc.Nonce(nonce), oauth2.AccessTypeOffline)
	s.log.Infof("OAuth redirect url: %s", url)
	return OIDCRequest{State: state, Nonce: nonce, AuthURL: url}, nil
}

func (s *service) OAuthRedirect(ctx context.Context, authProvider, authType, code, state string, cookies []*http.Cookie) (*BookshelfTokens, apierr.Error) {
	stateCookie := request.FilterCookies(cookies, "state")
	if stateCookie == nil {
		s.log.Error("no state cookie in OAuth redirect")
		return nil, apierr.NewBadRequestError("no cookies in auth request")
	}
	if state != stateCookie.Value {
		s.log.Error("state values did not match: %s - %s")
		return nil, apierr.NewBadRequestError("invalid token")
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
		return nil, apierr.NewBadRequestError("invalid auth provider in request url")
	}
	tokens, err := NewTokens(s.log, APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return nil, apierr.NewInternalServerError()
	}
	err = s.db.NewRefreshToken(ctx, APIKey, tokens.refreshToken)
	if err != nil {
		s.log.Error("could not save refresh token to db")
		return nil, apierr.NewInternalServerError()
	}
	return tokens, nil
}

func (s *service) RefreshTokens(ctx context.Context, accessToken, code string) (*BookshelfTokens, apierr.Error) {
	tkn, err := ParseJWT(s.log, accessToken)
	if err != nil || !tkn.HasCorrectClaims(code) {
		s.log.Error("could not parse jwt from cookie")
		return nil, apierr.NewJWTTokenError("could not parse token")
	}
	APIKey := tkn.Subject
	token, err := s.db.GetRefreshTokenByAPIKey(ctx, APIKey)
	if err != nil {
		s.log.Error("could not get refresh token from db")
		if err == apierr.ErrInternalServerError {
			return nil, apierr.NewInternalServerError()
		}
		return nil, apierr.NewAPIError(http.StatusNotFound, err, "no refresh token")
	}
	tkn, err = ParseJWT(s.log, token)
	if err != nil || !tkn.IsValid() || !tkn.HasCorrectClaims(code) {
		s.log.Error("parsed refresh token invalid")
		return nil, apierr.NewBadRequestError("invalid refresh token")
	}
	tokens, err := NewTokens(s.log, APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return nil, apierr.NewInternalServerError()
	}
	err = s.db.NewRefreshToken(ctx, APIKey, tokens.refreshToken)
	if err != nil {
		s.log.Error("could not save refresh token to db")
		return nil, apierr.NewInternalServerError()
	}
	return tokens, nil
}

func (s *service) LogOut(ctx context.Context, APIKey string) apierr.Error {
	numDeleted, err := s.db.DeleteRefreshToken(ctx, APIKey)
	if err != nil {
		s.log.Error("error deleting refresh token from db: %+v", err)
		return apierr.NewInternalServerError()
	}
	if numDeleted != 1 {
		s.log.Error("attempted to delete one refresh token, got: %d", numDeleted)
		return apierr.NewBadRequestError("incorrect APIKey")
	}
	return nil
}
