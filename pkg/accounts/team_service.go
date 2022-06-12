package accounts

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
)

// TeamRepository provides access to team storage.
type TeamRepository interface {
	New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr)
	DeleteTeam(ctx context.Context, requestData DelTeamRequest) (int, errors.ApiErr)
	AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr)
	DelSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr)
	DelMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr)
	AddCmdToTeam(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr)
	DelCmdFromTeam(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr)
}

// TeamService provides Team operations.
type TeamService interface {
	New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr)
	DeleteTeam(ctx context.Context, requestData DelTeamRequest) (int, errors.ApiErr)
	AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr)
	DelSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr)
	DelMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr)
	AddCmdToTeam(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr)
	DelCmdFromTeam(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr)
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
	teamID, err := s.r.New(ctx, requestData)
	return teamID, err
}

func (s *teamService) DeleteTeam(ctx context.Context, requestData DelTeamRequest) (int, errors.ApiErr) {
	numDeleted, err := s.r.DeleteTeam(ctx, requestData)
	return numDeleted, err
}

// AddMember calls the repository method for adding a member to a new team.
func (s *teamService) AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr) {
	ok, err := s.r.AddMember(ctx, requestData)
	return ok, err
}

func (s *teamService) DelSelf(ctx context.Context, requestData DelSelfRequest) (bool, errors.ApiErr) {
	ok, err := s.r.DelSelf(ctx, requestData)
	return ok, err
}

func (s *teamService) DelMember(ctx context.Context, requestData DelMemberRequest) (bool, errors.ApiErr) {
	if requestData.Role == RoleUser {
		return false, errors.NewWrongCredentialsError("user not authorized to remove members from teams")
	}
	ok, err := s.r.DelMember(ctx, requestData)
	return ok, err
}

// AddCmdToTeam calls the repository method for adding a command to a teams bookmarks.
func (s *teamService) AddCmdToTeam(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr) {
	numUpdated, err := s.r.AddCmdToTeam(ctx, requestData)
	return numUpdated, err
}

// DelCmdFromTeam calls the repository method for removing a command from a teams bookmarks.
func (s *teamService) DelCmdFromTeam(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr) {
	numUpdated, err := s.r.DelCmdFromTeam(ctx, requestData, APIKey)
	return numUpdated, err
}
