package accounts

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
)

// UserRepository provides access to the user storage.
type UserRepository interface {
	NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.ApiErr)
	LogIn(ctx context.Context, requestData LogInRequest) (User, errors.ApiErr)
	// GetTeams(ctx context.Context, APIKey string) ([]Team, errors.ApiErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.ApiErr)
	AddCmd(reqCtx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.ApiErr)
	DeleteCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.ApiErr)
	Delete(reqCtx context.Context, requestData DelUserRequest, APIKey string) (int, errors.ApiErr)
}

// UserService provides the user operations.
type UserService interface {
	NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.ApiErr)
	LogIn(ctx context.Context, requestData LogInRequest) (User, errors.ApiErr)
	// GetTeams(ctx context.Context, APIKey string) ([]Team, errors.ApiErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.ApiErr)
	AddCmd(reqCtx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.ApiErr)
	DeleteCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.ApiErr)
	Delete(ctx context.Context, requestData DelUserRequest, APIKey string) (int, errors.ApiErr)
}

type userService struct {
	r UserRepository
}

// NewUserService creates a search service with the necessary dependencies.
func NewUserService(r UserRepository) UserService {
	return &userService{r}
}

// Search returns the url of a given cmd.
func (s *userService) NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.ApiErr) {
	// TODO: add validation here
	user, err := s.r.NewUser(ctx, requestData)
	return user, err
}

// Login takes in request data, checks the db and returns the username and apikey is successful.
func (s *userService) LogIn(ctx context.Context, requestData LogInRequest) (User, errors.ApiErr) {
	currUser, err := s.r.LogIn(ctx, requestData)
	return currUser, err
}

// // GetTeams calls the GetTeams method and returns all teams for a user.
// func (s *userService) GetTeams(ctx context.Context, APIKey string) ([]Team, errors.ApiErr) {
// 	teams, err := s.r.GetTeams(ctx, APIKey)
// 	return teams, err
// }

// GetAllCmds calls the GetAllCmds method and returns all the users commands.
func (s *userService) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.ApiErr) {
	cmds, err := s.r.GetAllCmds(ctx, APIKey)
	return cmds, err
}

// AddCmd calls the AddCmd method and returns the number of updated commands.
func (s *userService) AddCmd(ctx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.ApiErr) {
	numUpdated, err := s.r.AddCmd(ctx, requestData, APIKey)
	return numUpdated, err
}

// DelCmd calls the DelCmd method and returns the number of updated commands.
func (s *userService) DeleteCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.ApiErr) {
	numUpdated, err := s.r.DeleteCmd(ctx, requestData, APIKey)
	return numUpdated, err
}

// Delete calls the Delete method and returns the number of deleted users.
func (s *userService) Delete(ctx context.Context, requestData DelUserRequest, APIKey string) (int, errors.ApiErr) {
	// TODO: add validation here
	user, err := s.r.Delete(ctx, requestData, APIKey)
	return user, err
}
