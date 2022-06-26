package accounts

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
)

// TeamRepository provides access to team storage.
type TeamRepository interface {
	NewTeam(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr)
	DeleteTeam(ctx context.Context, requestData DelTeamRequest) (int, errors.ApiErr)
	DeleteSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr)
	AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr)
	DeleteMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr)
	AddTeamCmd(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr)
	DeleteTeamCmd(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr)
}

// TeamService provides Team operations.
type TeamService interface {
	New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr)
	Delete(ctx context.Context, requestData DelTeamRequest) (int, errors.ApiErr)
	AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr)
	DeleteSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr)
	DeleteMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr)
	AddCmd(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr)
	DeleteCmd(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr)
}

type teamService struct {
	r TeamRepository
}

// NewTeamService creates a search service with the necessary dependencies.
func NewTeamService(r TeamRepository) TeamService {
	return &teamService{r}
}

// New calls the repository method for creating a new team.
func (s *teamService) New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr) {
	teamID, err := s.r.NewTeam(ctx, requestData)
	return teamID, err
}

func (s *teamService) Delete(ctx context.Context, requestData DelTeamRequest) (int, errors.ApiErr) {
	numDeleted, err := s.r.DeleteTeam(ctx, requestData)
	return numDeleted, err
}

// AddMember calls the repository method for adding a member to a new team.
func (s *teamService) AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr) {
	ok, err := s.r.AddMember(ctx, requestData)
	return ok, err
}

func (s *teamService) DeleteSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr) {
	ok, err := s.r.DeleteSelf(ctx, requestData)
	return ok, err
}

func (s *teamService) DeleteMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr) {
	if requestData.Role == RoleUser {
		return false, errors.NewWrongCredentialsError("user not authorized to remove members from teams")
	}
	ok, err := s.r.DeleteMember(ctx, requestData)
	return ok, err
}

// AddCmd calls the repository method for adding a command to a teams bookmarks.
func (s *teamService) AddCmd(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr) {
	numUpdated, err := s.r.AddTeamCmd(ctx, requestData)
	return numUpdated, err
}

// DeleteCmd calls the repository method for removing a command from a teams bookmarks.
func (s *teamService) DeleteCmd(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr) {
	numUpdated, err := s.r.DeleteTeamCmd(ctx, requestData, APIKey)
	return numUpdated, err
}
