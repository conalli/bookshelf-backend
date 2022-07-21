package search

import (
	"context"
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
	GetAllBookmarks(ctx context.Context, APIKey string) ([]accounts.Bookmark, errors.APIErr)
	AddBookmark(reqCtx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr)
}

// Service provides the search operation.
type Service interface {
	Search(ctx context.Context, APIKey, cmd string) (any, error)
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
func (s *service) Search(ctx context.Context, APIKey, cmd string) (any, error) {
	ctx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Var(APIKey, "uuid")
	if err != nil {
		s.log.Error("invalid API key")
		return "", errors.NewBadRequestError("invalid API key")
	}
	cmds := strings.Fields(cmd)
	switch cmds[0] {
	case "ls":
		ls := NewLSFlagset()
		err := ls.fs.Parse(cmds[1:])
		if err != nil || *ls.b && *ls.c {
			s.log.Error("could not parse ls flag cmds")
			return "", errors.NewBadRequestError("bad ls flags")
		}
		if *ls.b {
			return s.db.GetAllBookmarks(ctx, APIKey)
		}
		if *ls.c {
			usr, err := s.db.GetUserByAPIKey(ctx, APIKey)
			return usr.Cmds, err
		}
	case "touch", "add":
		touch := NewTouchFlagset()
		err := touch.fs.Parse(cmds[1:])
		if err != nil {
			s.log.Error("could not parse ls flag cmds")
			return "", errors.NewBadRequestError("bad ls flags")
		}
		if *touch.b && len(*touch.url) > 0 {
			req := request.AddBookmark{
				Name: *touch.name,
				URL:  *touch.url,
				Path: *touch.path,
			}
			return s.db.AddBookmark(ctx, req, APIKey)
		}
	default:
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
	return nil, nil
}

func (s *service) Evaluate(args []string) error {
	return nil
}
