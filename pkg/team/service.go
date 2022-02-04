package team

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
)

// Repository provides access to team storage.
type Repository interface {
	New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr)
	AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr)
	AddCmdToTeam(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr)
	DelCmdFromTeam(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr)
}

// Service provides reviewing operations.
type Service interface {
	New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr)
	AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr)
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

func (s *service) New(ctx context.Context, requestData NewTeamRequest) (string, errors.ApiErr) {
	teamID, err := s.r.New(ctx, requestData)
	return teamID, err
}

func (s *service) AddMember(ctx context.Context, requestData AddMemberRequest) (bool, errors.ApiErr) {
	ok, err := s.r.AddMember(ctx, requestData)
	return ok, err
}

func (s *service) AddCmdToTeam(ctx context.Context, requestData AddTeamCmdRequest) (int, errors.ApiErr) {
	numUpdated, err := s.r.AddCmdToTeam(ctx, requestData)
	return numUpdated, err
}

func (s *service) DelCmdFromTeam(ctx context.Context, requestData DelTeamCmdRequest, APIKey string) (int, errors.ApiErr) {
	numUpdated, err := s.r.DelCmdFromTeam(ctx, requestData, APIKey)
	return numUpdated, err
}
