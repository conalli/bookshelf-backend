package user

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
)

// Repository provides access to the user storage.
type Repository interface {
	NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.ApiErr)
	LogIn(ctx context.Context, requestData LogInRequest) (User, errors.ApiErr)
	GetTeams(ctx context.Context, APIKey string) ([]Team, errors.ApiErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.ApiErr)
	AddCmd(reqCtx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.ApiErr)
	DelCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.ApiErr)
	Delete(reqCtx context.Context, requestData DelUserRequest, APIKey string) (int, errors.ApiErr)
}

// Service provides the user operations.
type Service interface {
	NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.ApiErr)
	LogIn(ctx context.Context, requestData LogInRequest) (User, errors.ApiErr)
	GetTeams(ctx context.Context, APIKey string) ([]Team, errors.ApiErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.ApiErr)
	AddCmd(reqCtx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.ApiErr)
	DelCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.ApiErr)
	Delete(ctx context.Context, requestData DelUserRequest, APIKey string) (int, errors.ApiErr)
}

type service struct {
	r Repository
}

// NewService creates a search service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// Search returns the url of a given cmd.
func (s *service) NewUser(ctx context.Context, requestData SignUpRequest) (User, errors.ApiErr) {
	// TODO: add validation here
	user, err := s.r.NewUser(ctx, requestData)
	return user, err
}

// Login takes in request data, checks the db and returns the username and apikey is successful.
func (s *service) LogIn(ctx context.Context, requestData LogInRequest) (User, errors.ApiErr) {
	currUser, err := s.r.LogIn(ctx, requestData)
	return currUser, err
}

// GetTeams calls the GetTeams method and returns all teams for a user.
func (s *service) GetTeams(ctx context.Context, APIKey string) ([]Team, errors.ApiErr) {
	teams, err := s.r.GetTeams(ctx, APIKey)
	return teams, err
}

// GetAllCmds calls the GetAllCmds method and returns all the users commands.
func (s *service) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.ApiErr) {
	cmds, err := s.r.GetAllCmds(ctx, APIKey)
	return cmds, err
}

// AddCmd calls the AddCmd method and returns the number of updated commands.
func (s *service) AddCmd(ctx context.Context, requestData AddCmdRequest, APIKey string) (int, errors.ApiErr) {
	numUpdated, err := s.r.AddCmd(ctx, requestData, APIKey)
	return numUpdated, err
}

// DelCmd calls the DelCmd method and returns the number of updated commands.
func (s *service) DelCmd(ctx context.Context, requestData DelCmdRequest, APIKey string) (int, errors.ApiErr) {
	numUpdated, err := s.r.DelCmd(ctx, requestData, APIKey)
	return numUpdated, err
}

// Delete calls the Delete method and returns the number of deleted users.
func (s *service) Delete(ctx context.Context, requestData DelUserRequest, APIKey string) (int, errors.ApiErr) {
	// TODO: add validation here
	user, err := s.r.Delete(ctx, requestData, APIKey)
	return user, err
}
