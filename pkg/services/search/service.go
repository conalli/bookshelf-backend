package search

import (
	"context"
	"fmt"

	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

// Repository provides access to storage.
type Repository interface {
	GetUserByAPIKey(ctx context.Context, APIKey string) (accounts.User, error)
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
	ctx, cancelFunc := request.WithDefaultTimeout(ctx)
	defer cancelFunc()
	// TODO: add validation here for team/ user cmd
	usr, err := s.r.GetUserByAPIKey(ctx, APIKey)
	defaultSearch := fmt.Sprintf("http://www.google.com/search?q=%s", cmd)
	if err != nil {
		return defaultSearch, err
	}
	url, ok := usr.Bookmarks[cmd]
	if !ok {
		return defaultSearch, err
	}
	return formatURL(url), err
}
