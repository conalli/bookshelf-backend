package team

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
)

// Repository provides access to team storage.
type Repository interface {
	New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr)
	AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr)
	DelSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr)
	DelMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr)
	AddCmdToTeam(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr)
	DelCmdFromTeam(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr)
}

// Service provides Team operations.
type Service interface {
	New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr)
	AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr)
	DelSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr)
	DelMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr)
	AddCmdToTeam(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr)
	DelCmdFromTeam(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr)
}

type service struct {
	r Repository
}

// NewService creates a search service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// New calls the repository method for creating a new team.
func (s *service) New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr) {
	teamID, err := s.r.New(ctx, requestData)
	return teamID, err
}

// AddMember calls the repository method for adding a member to a new team.
func (s *service) AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr) {
	ok, err := s.r.AddMember(ctx, requestData)
	return ok, err
}

func (s *service) DelSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr) {
	ok, err := s.r.DelSelf(ctx, requestData)
	return ok, err
}

func (s *service) DelMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr) {
	if requestData.Role == RoleUser {
		return false, errors.NewWrongCredentialsError("user not authorized to remove members from teams")
	}
	ok, err := s.r.DelMember(ctx, requestData)
	return ok, err
}

// AddCmdToTeam calls the repository method for adding a command to a teams bookmarks.
func (s *service) AddCmdToTeam(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr) {
	numUpdated, err := s.r.AddCmdToTeam(ctx, requestData)
	return numUpdated, err
}

// DelCmdFromTeam calls the repository method for removing a command from a teams bookmarks.
func (s *service) DelCmdFromTeam(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr) {
	numUpdated, err := s.r.DelCmdFromTeam(ctx, requestData, APIKey)
	return numUpdated, err
}
