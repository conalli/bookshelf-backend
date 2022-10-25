package bookmarks

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	GetAllBookmarks(ctx context.Context, APIKey string) ([]Bookmark, errors.APIErr)
	GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]Bookmark, errors.APIErr)
	AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr)
	DeleteBookmark(ctx context.Context, requestData request.DeleteBookmark, APIKey string) (int, errors.APIErr)
}

type Repository interface {
	GetAllBookmarks(ctx context.Context, APIKey string) ([]Bookmark, errors.APIErr)
	GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]Bookmark, errors.APIErr)
	AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr)
	DeleteBookmark(ctx context.Context, requestData request.DeleteBookmark, APIKey string) (int, errors.APIErr)
}

type service struct {
	log      logs.Logger
	validate *validator.Validate
	db       Repository
}

func NewService(l logs.Logger, v *validator.Validate, db Repository) *service {
	return &service{l, v, db}
}

func (s *service) GetAllBookmarks(ctx context.Context, APIKey string) ([]Bookmark, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Var(APIKey, "uuid")
	if validateErr != nil {
		s.log.Errorf("Could not validate GET ALL BOOKMARKS request: %v", validateErr)
		return nil, errors.NewBadRequestError("request format incorrect.")
	}
	books, err := s.db.GetAllBookmarks(reqCtx, APIKey)
	return books, err
}

func (s *service) GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]Bookmark, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Var(APIKey, "uuid")
	if validateErr != nil {
		s.log.Errorf("Could not validate GET BOOKMARKS FOLDER request: %v", validateErr)
		return nil, errors.NewBadRequestError("request format incorrect.")
	}
	books, err := s.db.GetBookmarksFolder(reqCtx, path, APIKey)
	return books, err
}

// AddBookmark adds a bookmark for an account.
func (s *service) AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate ADD BOOKMARK request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.AddBookmark(reqCtx, requestData, APIKey)
	return numUpdated, err
}

// DeleteBookmark removes a bookmark from an account.
func (s *service) DeleteBookmark(ctx context.Context, requestData request.DeleteBookmark, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate DELETE BOOKMARK request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.DeleteBookmark(reqCtx, requestData, APIKey)
	return numUpdated, err
}
