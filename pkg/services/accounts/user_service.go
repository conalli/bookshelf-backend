package accounts

import (
	"context"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/password"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserRepository provides access to the user storage.
type UserRepository interface {
	NewUser(ctx context.Context, requestData request.SignUp) (User, errors.APIErr)
	GetUserByName(ctx context.Context, requestData request.LogIn) (User, error)
	// GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr)
	AddCmd(reqCtx context.Context, requestData request.AddCmd, APIKey string) (int, errors.APIErr)
	DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, errors.APIErr)
	Delete(reqCtx context.Context, requestData request.DeleteUser, APIKey string) (int, errors.APIErr)
}

// UserService provides the user operations.
type UserService interface {
	NewUser(ctx context.Context, requestData request.SignUp) (User, errors.APIErr)
	LogIn(ctx context.Context, requestData request.LogIn) (User, errors.APIErr)
	// GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr)
	AddCmd(reqCtx context.Context, requestData request.AddCmd, APIKey string) (int, errors.APIErr)
	DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, errors.APIErr)
	Delete(ctx context.Context, requestData request.DeleteUser, APIKey string) (int, errors.APIErr)
}

type userService struct {
	log      *zap.SugaredLogger
	validate *validator.Validate
	db       UserRepository
}

// NewUserService creates a search service with the necessary dependencies.
func NewUserService(l *zap.SugaredLogger, v *validator.Validate, r UserRepository) UserService {
	return &userService{l, v, r}
}

// Search returns the url of a given cmd.
func (s *userService) NewUser(ctx context.Context, requestData request.SignUp) (User, errors.APIErr) {
	reqCtx, cancelFunc := request.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Struct(requestData)
	if validateErr != nil {
		s.log.Errorf("Could not validate SIGN UP request: %v", validateErr)
		return User{}, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.NewUser(reqCtx, requestData)
	return user, err
}

// Login takes in request data, checks the db and returns the username and apikey is successful.
func (s *userService) LogIn(ctx context.Context, requestData request.LogIn) (User, errors.APIErr) {
	reqCtx, cancelFunc := request.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Struct(requestData)
	if validateErr != nil {
		s.log.Errorf("Could not validate LOG IN request: %v", validateErr)
		return User{}, errors.NewBadRequestError("request format incorrect.")
	}
	usr, err := s.db.GetUserByName(reqCtx, requestData)
	if err != nil || !password.CheckHashedPassword(usr.Password, requestData.Password) {
		log.Printf("login getuserbykey %+v", err)
		return User{}, errors.NewAPIError(http.StatusUnauthorized, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	return usr, nil
}

// // GetTeams calls the GetTeams method and returns all teams for a user.
// func (s *userService) GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr) {
// reqCtx, cancelFunc := request.WithDefaultTimeout(ctx)
// defer cancelFunc()
// 	teams, err := s.db.GetTeams(reqCtx, APIKey)
// 	return teams, err
// }

// GetAllCmds calls the GetAllCmds method and returns all the users commands.
func (s *userService) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr) {
	reqCtx, cancelFunc := request.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Var(APIKey, "uuid")
	if validateErr != nil {
		s.log.Errorf("Could not validate GET ALL CMDS request: %v", validateErr)
		return nil, errors.NewBadRequestError("request format incorrect.")
	}
	cmds, err := s.db.GetAllCmds(reqCtx, APIKey)
	return cmds, err
}

// AddCmd calls the AddCmd method and returns the number of updated commands.
func (s *userService) AddCmd(ctx context.Context, requestData request.AddCmd, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate ADD CMD request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.AddCmd(reqCtx, requestData, APIKey)
	return numUpdated, err
}

// DelCmd calls the DelCmd method and returns the number of updated commands.
func (s *userService) DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate DELETE CMD request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.DeleteCmd(reqCtx, requestData, APIKey)
	return numUpdated, err
}

// Delete calls the Delete method and returns the number of deleted users.
func (s *userService) Delete(ctx context.Context, requestData request.DeleteUser, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.WithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate DELETE USER request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.Delete(reqCtx, requestData, APIKey)
	return user, err
}
