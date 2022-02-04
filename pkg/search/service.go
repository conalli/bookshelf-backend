package search

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/db"
)

// Repository provides access to the search storage.
type Repository interface {
	Search(ctx context.Context, APIKey, cmd string) (string, error)
}

// Service provides the search operation.
type Service interface {
	Search(ctx context.Context, APIKey, cmd string) (string, error)
}

type service struct {
	r Repository
}

// NewService creates a search service with the necessary dependencies.
func NewService(r Repository) Service {
	return &service{r}
}

// Search returns the url of a given cmd.
func (s *service) Search(ctx context.Context, APIKey, cmd string) (string, error) {
	ctx, cancelFunc := db.ReqContextWithTimeout(ctx)
	defer cancelFunc()
	// TODO: add validation here for team/ user cmd
	url, err := s.r.Search(ctx, APIKey, cmd)
	return formatURL(url), err
}
