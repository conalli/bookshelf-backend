package search

import (
	"context"
	"fmt"
	"os"
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
	AddBookmark(reqCtx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr)
}

// Cache provides access to Caching for the Search service.
type Cache interface {
	GetSearchData(ctx context.Context, APIKey, cmd string) (string, error)
	AddCmds(ctx context.Context, APIKey string, cmds map[string]string) bool
	DeleteCmds(ctx context.Context, APIKey string) bool
}

// Service provides the search operation.
type Service interface {
	Search(ctx context.Context, APIKey, cmd string) (string, error)
}

type service struct {
	log      logs.Logger
	validate *validator.Validate
	db       Repository
	cache    Cache
}

// NewService creates a search service with the necessary dependencies.
func NewService(l logs.Logger, v *validator.Validate, r Repository, c Cache) Service {
	return &service{l, v, r, c}
}

// Search returns the url of a given cmd.
func (s *service) Search(ctx context.Context, APIKey, args string) (string, error) {
	ctx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Var(APIKey, "uuid")
	if err != nil {
		s.log.Error("invalid API key")
		return "", errors.NewBadRequestError("invalid API key")
	}
	cmds := strings.Fields(args)
	return s.evaluateArgs(ctx, APIKey, cmds)
}

func (s *service) evaluateArgs(ctx context.Context, APIKey string, args []string) (string, error) {
	switch args[0] {
	case "help":
		s.log.Info("webcli: help")
		return fmt.Sprintf("%s/webcli/help", os.Getenv("ALLOWED_URL_BASE")), nil

	case "ls":
		ls := NewLSFlagset()
		err := ls.fs.Parse(args[1:])
		if err != nil || *ls.b && *ls.c {
			s.log.Error("webcli: could not parse ls flag cmds")
			return "", errors.NewBadRequestError("bad ls flags")
		}
		if *ls.b && *ls.c || len(*ls.bf) > 0 && *ls.c || *ls.b && len(*ls.bf) > 0 {
			s.log.Error("webcli: incorrect flags passed")
			return fmt.Sprintf("%s/404", os.Getenv("ALLOWED_URL_BASE")), nil
		}
		if *ls.b {
			s.log.Info("webcli: list bookmarks")
			return fmt.Sprintf("%s/webcli/bookmark?APIKey=%s", os.Getenv("ALLOWED_URL_BASE"), APIKey), nil
		}
		if *ls.c {
			s.log.Info("webcli: list commands")
			return fmt.Sprintf("%s/webcli/command?APIKey=%s", os.Getenv("ALLOWED_URL_BASE"), APIKey), nil
		}
		if ls.bf != nil {
			s.log.Info("webcli: list bookmark folder")
			return fmt.Sprintf("%s/webcli/bookmark?APIKey=%s&folder=%s", os.Getenv("ALLOWED_URL_BASE"), APIKey, *ls.bf), nil
		}
	case "touch", "add":
		touch := NewTouchFlagset()
		err := touch.fs.Parse(args[1:])
		if err != nil {
			s.log.Error("could not parse touch flag cmds")
			return "", errors.NewBadRequestError("bad touch flags")
		}
		if len(*touch.url) < 5 || *touch.b && len(*touch.c) > 0 {
			s.log.Error("webcli: incorrect flags passed")
			return fmt.Sprintf("%s/404", os.Getenv("ALLOWED_URL_BASE")), nil
		}
		if *touch.b {
			req := request.AddBookmark{
				Name: *touch.name,
				URL:  *touch.url,
				Path: *touch.path,
			}
			res, err := s.db.AddBookmark(ctx, req, APIKey)
			if err != nil {
				return "", err
			}
			if res == 0 {
				return fmt.Sprintf("%s/404", os.Getenv("ALLOWED_URL_BASE")), nil
			}
			return fmt.Sprintf("%s/webcli/success", os.Getenv("ALLOWED_URL_BASE")), nil
		}
	default:
		cachedURL, err := s.cache.GetSearchData(ctx, APIKey, args[0])
		if err == nil {
			s.log.Info("retrieved search data from cache")
			return formatURL(cachedURL), nil
		}
		s.log.Infof("could not get search data from cache: %v", err)
		defaultSearch := fmt.Sprintf("http://www.google.com/search?q=%s", args[0])
		usr, err := s.db.GetUserByAPIKey(ctx, APIKey)
		if err != nil {
			s.log.Errorf("could not get user by API key: %v", err)
			return defaultSearch, err
		}
		ok := s.cache.AddCmds(ctx, APIKey, usr.Cmds)
		if !ok {
			s.log.Errorf("could not add cmds to cache: %v", err)
		}
		url, ok := usr.Cmds[args[0]]
		if !ok {
			s.log.Infof("Cmd %s does not exist. Returning default search", args[0])
			return defaultSearch, nil
		}
		return formatURL(url), nil
	}
	return "", nil
}
