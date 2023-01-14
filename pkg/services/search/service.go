package search

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/go-playground/validator/v10"
)

// Repository provides access to storage.
type Repository interface {
	GetUserByAPIKey(ctx context.Context, APIKey string) (accounts.User, error)
	AddBookmark(reqCtx context.Context, requestData request.AddBookmark, APIKey string) (int, apierr.Error)
	NewRefreshToken(ctx context.Context, APIKey, refreshToken string) error
	GetRefreshTokenByAPIKey(ctx context.Context, APIKey string) (string, error)
}

// Cache provides access to Caching for the Search service.
type Cache interface {
	GetAllCmds(ctx context.Context, cacheKey string) (map[string]string, error)
	GetOneCmd(ctx context.Context, cacheKey, cmd string) (string, error)
	AddCmds(ctx context.Context, cacheKey string, cmds map[string]string) (int64, error)
	DeleteCmds(ctx context.Context, cacheKey string) (int64, error)
}

// Service provides the search operation.
type Service interface {
	Search(ctx context.Context, APIKey, args, code string, refresh bool) (string, *auth.BookshelfTokens, error)
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

type refreshResult struct {
	tkn *auth.BookshelfTokens
	err error
}

// Search returns the url of a given cmd.
func (s *service) Search(ctx context.Context, APIKey, args, code string, refresh bool) (string, *auth.BookshelfTokens, error) {
	ctx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	err := s.validate.Var(APIKey, "uuid")
	if err != nil {
		s.log.Error("invalid API key")
		return "", nil, apierr.NewBadRequestError("invalid API key")
	}
	refChan := make(chan refreshResult, 1)
	if refresh {
		go func() {
			tokens, err := s.refresh(ctx, APIKey, code)
			res := refreshResult{tokens, err}
			refChan <- res
			close(refChan)
		}()
	} else {
		refChan <- refreshResult{}
		close(refChan)
	}
	cmds := strings.Fields(args)
	url, err := s.evaluateArgs(ctx, APIKey, cmds)
	if err != nil {
		s.log.Error("could not evaluate args in search")
		return "", nil, err
	}
	res := <-refChan
	if res.err != nil {
		s.log.Error("could not refresh tokens in search")
		return "", nil, err
	}
	return url, res.tkn, nil
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
			return "", apierr.NewBadRequestError("bad ls flags")
		}
		if *ls.b && *ls.c || len(*ls.bf) > 0 && *ls.c || *ls.b && len(*ls.bf) > 0 {
			s.log.Error("webcli: incorrect flags passed")
			return fmt.Sprintf("%s/404", os.Getenv("ALLOWED_URL_BASE")), nil
		}
		if *ls.b {
			s.log.Info("webcli: list bookmarks")
			return fmt.Sprintf("%s/webcli/bookmark", os.Getenv("ALLOWED_URL_BASE")), nil
		}
		if *ls.bf != "" {
			s.log.Infof("FLAG: %s", *ls.bf)
			s.log.Info("webcli: list bookmark folder")
			return fmt.Sprintf("%s/webcli/bookmark?folder=%s", os.Getenv("ALLOWED_URL_BASE"), *ls.bf), nil
		}
		if *ls.c {
			s.log.Info("webcli: list commands")
			return fmt.Sprintf("%s/webcli/command", os.Getenv("ALLOWED_URL_BASE")), nil
		}
	case "touch", "add":
		touch := NewTouchFlagset()
		err := touch.fs.Parse(args[1:])
		if err != nil {
			s.log.Error("could not parse touch flag cmds")
			return "", apierr.NewBadRequestError("bad touch flags")
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
		cachedURL, err := s.cache.GetOneCmd(ctx, APIKey, args[0])
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
		numAdded, err := s.cache.AddCmds(ctx, APIKey, usr.Cmds)
		if err != nil {
			s.log.Errorf("could not add cmds to cache: %v", err)
		}
		if numAdded == 0 {
			s.log.Error("could not add cmds to cache")
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

func (s *service) refresh(ctx context.Context, APIKey, code string) (*auth.BookshelfTokens, error) {
	token, err := s.db.GetRefreshTokenByAPIKey(ctx, APIKey)
	if err != nil {
		s.log.Error("could not get refresh token from db")
		if err == apierr.ErrInternalServerError {
			return nil, apierr.ErrInternalServerError
		}
		return nil, apierr.ErrNotFound
	}
	tkn, err := auth.ParseJWT(s.log, token)
	if err != nil || !tkn.IsValid() || !tkn.HasCorrectClaims(code) {
		s.log.Error("parsed refresh token invalid")
		return nil, apierr.ErrBadRequest
	}
	tokens, err := auth.NewTokens(s.log, APIKey)
	if err != nil {
		s.log.Error("could not create new tokens")
		return nil, apierr.ErrInternalServerError
	}
	err = s.db.NewRefreshToken(ctx, APIKey, tokens.RefreshToken())
	if err != nil {
		s.log.Error("could not save refresh token to db")
		return nil, apierr.ErrInternalServerError
	}
	return tokens, nil
}
