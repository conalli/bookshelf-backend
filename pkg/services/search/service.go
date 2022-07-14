package search

import (
	"context"
	"fmt"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/go-playground/validator/v10"
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
	log      logs.Logger
	validate *validator.Validate
	db       Repository
}

// NewService creates a search service with the necessary dependencies.
func NewService(l logs.Logger, v *validator.Validate, r Repository) Service {
	return &service{l, v, r}
}

// Search returns the url of a given cmd.
func (s *service) Search(ctx context.Context, APIKey, cmd string) (string, error) {
	ctx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Var(APIKey, "uuid")
	if err != nil {
		return "", errors.NewBadRequestError("invalid API key")
	}
	usr, err := s.db.GetUserByAPIKey(ctx, APIKey)
	defaultSearch := fmt.Sprintf("http://www.google.com/search?q=%s", cmd)
	if err != nil {
		s.log.Errorf("Could not GET USER BY API KEY: %v", err)
		return defaultSearch, err
	}
	url, ok := usr.Bookmarks[cmd]
	if !ok {
		s.log.Infof("Cmd %s does not exist. Returning default search", cmd)
		return defaultSearch, nil
	}
	return formatURL(url), nil
}
