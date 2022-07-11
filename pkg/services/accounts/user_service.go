package accounts

import (
	"context"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/reqcontext"
	"github.com/conalli/bookshelf-backend/pkg/password"
	"github.com/go-playground/validator/v10"
)

// UserRepository provides access to the user storage.
type UserRepository interface {
	NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.APIErr)
	GetUserByName(ctx context.Context, requestData LogInRequest) (User, error)
	// GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr)
	AddCmd(reqCtx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.APIErr)
	DeleteCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.APIErr)
	Delete(reqCtx context.Context, requestData DelUserRequest, APIKey string) (int, errors.APIErr)
}

// UserService provides the user operations.
type UserService interface {
	NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.APIErr)
	LogIn(ctx context.Context, requestData LogInRequest) (User, errors.APIErr)
	// GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr)
	AddCmd(reqCtx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.APIErr)
	DeleteCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.APIErr)
	Delete(ctx context.Context, requestData DelUserRequest, APIKey string) (int, errors.APIErr)
}

type userService struct {
	v *validator.Validate
	r UserRepository
}

// NewUserService creates a search service with the necessary dependencies.
func NewUserService(v *validator.Validate, r UserRepository) UserService {
	return &userService{v, r}
}

// Search returns the url of a given cmd.
func (s *userService) NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.APIErr) {
	reqCtx, cancelFunc := reqcontext.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.v.Struct(requestData)
	if validateErr != nil {
		return User{}, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.r.NewUser(reqCtx, requestData)
	return user, err
}

// Login takes in request data, checks the db and returns the username and apikey is successful.
func (s *userService) LogIn(ctx context.Context, requestData LogInRequest) (User, errors.APIErr) {
	reqCtx, cancelFunc := reqcontext.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.v.Struct(requestData)
	if validateErr != nil {
		return User{}, errors.NewBadRequestError("request format incorrect.")
	}
	usr, err := s.r.GetUserByName(reqCtx, requestData)
	if err != nil || !password.CheckHashedPassword(usr.Password, requestData.Password) {
		log.Printf("login getuserbykey %+v", err)
		return User{}, errors.NewAPIError(http.StatusUnauthorized, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	return usr, nil
}

// // GetTeams calls the GetTeams method and returns all teams for a user.
// func (s *userService) GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr) {
// reqCtx, cancelFunc := reqcontext.WithDefaultTimeout(ctx)
// defer cancelFunc()
// 	teams, err := s.r.GetTeams(reqCtx, APIKey)
// 	return teams, err
// }

// GetAllCmds calls the GetAllCmds method and returns all the users commands.
func (s *userService) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr) {
	reqCtx, cancelFunc := reqcontext.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.v.Var(APIKey, "uuid")
	if validateErr != nil {
		return nil, errors.NewBadRequestError("request format incorrect.")
	}
	cmds, err := s.r.GetAllCmds(reqCtx, APIKey)
	return cmds, err
}

// AddCmd calls the AddCmd method and returns the number of updated commands.
func (s *userService) AddCmd(ctx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := reqcontext.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.v.Struct(requestData)
	validateAPIKeyErr := s.v.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.r.AddCmd(reqCtx, requestData, APIKey)
	return numUpdated, err
}

// DelCmd calls the DelCmd method and returns the number of updated commands.
func (s *userService) DeleteCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := reqcontext.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.v.Struct(requestData)
	validateAPIKeyErr := s.v.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.r.DeleteCmd(reqCtx, requestData, APIKey)
	return numUpdated, err
}

// Delete calls the Delete method and returns the number of deleted users.
func (s *userService) Delete(ctx context.Context, requestData DelUserRequest, APIKey string) (int, errors.APIErr) {
	// TODO: add validation here
	reqCtx, cancelFunc := reqcontext.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.v.Struct(requestData)
	validateAPIKeyErr := s.v.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.r.Delete(reqCtx, requestData, APIKey)
	return user, err
}
