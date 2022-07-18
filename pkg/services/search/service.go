package search

import (
	"context"
	"flag"
	"fmt"
	"strings"

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
	flags    map[string]flag.FlagSet
}

// NewService creates a search service with the necessary dependencies.
func NewService(l logs.Logger, v *validator.Validate, r Repository) Service {
	f := make(map[string]flag.FlagSet)
	return &service{l, v, r, f}
}

// Search returns the url of a given cmd.
func (s *service) Search(ctx context.Context, APIKey, cmd string) (string, error) {
	ctx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Var(APIKey, "uuid")
	if err != nil {
		s.log.Error("invalid API key")
		return "", errors.NewBadRequestError("invalid API key")
	}
	strings.Fields(cmd)
	usr, err := s.db.GetUserByAPIKey(ctx, APIKey)
	defaultSearch := fmt.Sprintf("http://www.google.com/search?q=%s", cmd)
	if err != nil {
		s.log.Errorf("could not get user by API key: %v", err)
		return defaultSearch, err
	}
	url, ok := usr.Cmds[cmd]
	if !ok {
		s.log.Infof("Cmd %s does not exist. Returning default search", cmd)
		return defaultSearch, nil
	}
	return formatURL(url), nil
}

func (s *service) Evaluate(args []string) error {
	return nil
}
