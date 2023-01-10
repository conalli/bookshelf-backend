package accounts

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/go-playground/validator/v10"
)

// UserRepository provides access to the user storage.
type UserRepository interface {
	GetUserByAPIKey(ctx context.Context, APIKey string) (User, error)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, apierr.Error)
	AddCmd(reqCtx context.Context, requestData request.AddCmd, APIKey string) (int, apierr.Error)
	DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, apierr.Error)
	Delete(reqCtx context.Context, requestData request.DeleteUser, APIKey string) (int, apierr.Error)
}

// UserCache provides access to the cache.
type UserCache interface {
	GetUser(ctx context.Context, userKey string) (User, error)
	DeleteUser(ctx context.Context, userKey string) (int64, error)
	GetAllCmds(ctx context.Context, cacheKey string) (map[string]string, error)
	AddCmds(ctx context.Context, cacheKey string, cmds map[string]string) (int64, error)
	DeleteCmds(ctx context.Context, cacheKey string) (int64, error)
}

// UserService provides the user operations.
type UserService interface {
	UserInfo(ctx context.Context, APIKey string) (User, apierr.Error)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, apierr.Error)
	AddCmd(reqCtx context.Context, requestData request.AddCmd, APIKey string) (int, apierr.Error)
	DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, apierr.Error)
	Delete(ctx context.Context, requestData request.DeleteUser, APIKey string) (int, apierr.Error)
}

type userService struct {
	log      logs.Logger
	validate *validator.Validate
	db       UserRepository
	cache    UserCache
}

// NewUserService creates a search service with the necessary dependencies.
func NewUserService(l logs.Logger, v *validator.Validate, r UserRepository, c UserCache) UserService {
	return &userService{l, v, r, c}
}

func (s *userService) UserInfo(ctx context.Context, APIKey string) (User, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Var(APIKey, "uuid")
	if validateErr != nil {
		s.log.Errorf("could not validate GET ALL CMDS request: %v", validateErr)
		return User{}, apierr.NewBadRequestError("request format incorrect.")
	}
	user, err := s.cache.GetUser(ctx, APIKey)
	if err != nil {
		s.log.Error("could not get user from cache")
	}
	if user.APIKey != "" {
		s.log.Info("got user from cache")
		return user, nil
	}
	user, err = s.db.GetUserByAPIKey(reqCtx, APIKey)
	if err != nil {
		s.log.Error("could not get user by APIKey: %v", err)
		return User{}, apierr.NewInternalServerError()
	}
	return user, nil
}

// Delete calls the Delete method and returns the number of deleted users.
func (s *userService) Delete(ctx context.Context, requestData request.DeleteUser, APIKey string) (int, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("could not validate DELETE USER request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, apierr.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.Delete(reqCtx, requestData, APIKey)
	s.cache.DeleteUser(ctx, APIKey)
	return user, err
}
