package bookmarks

import (
	"context"
	"net/http"
	"regexp"

	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	GetAllBookmarks(ctx context.Context, APIKey string) (*Folder, apierr.Error)
	GetBookmarksFolder(ctx context.Context, path, APIKey string) (*Folder, apierr.Error)
	AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, apierr.Error)
	AddBookmarksFromFile(ctx context.Context, r *http.Request, APIKey string) (int, apierr.Error)
	DeleteBookmark(ctx context.Context, bookmarkID, APIKey string) (int, apierr.Error)
}

type Repository interface {
	GetAllBookmarks(ctx context.Context, APIKey string) ([]Bookmark, apierr.Error)
	GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]Bookmark, apierr.Error)
	AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, apierr.Error)
	AddManyBookmarks(ctx context.Context, bookmarks []Bookmark) (int, apierr.Error)
	DeleteBookmark(ctx context.Context, bookmarkID, APIKey string) (int, apierr.Error)
}

type service struct {
	log      logs.Logger
	validate *validator.Validate
	db       Repository
}

func NewService(l logs.Logger, v *validator.Validate, db Repository) *service {
	return &service{l, v, db}
}

func (s *service) GetAllBookmarks(ctx context.Context, APIKey string) (*Folder, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Var(APIKey, "uuid")
	if validateErr != nil {
		s.log.Errorf("Could not validate GET ALL BOOKMARKS request: %v", validateErr)
		return nil, apierr.NewBadRequestError("request format incorrect.")
	}
	books, err := s.db.GetAllBookmarks(reqCtx, APIKey)
	folder := organizeBookmarks(books, "", BookmarksBasePath, BookmarksBasePath, BookmarksBasePath)
	return folder, err
}

func (s *service) GetBookmarksFolder(ctx context.Context, folderName, APIKey string) (*Folder, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Var(APIKey, "uuid")
	if validateErr != nil {
		s.log.Errorf("Could not validate GET BOOKMARKS FOLDER request: %v", validateErr)
		return nil, apierr.NewBadRequestError("request format incorrect.")
	}
	books, err := s.db.GetBookmarksFolder(reqCtx, folderName, APIKey)
	if err != nil {
		s.log.Error("could not get bookmarks from folder %s", folderName)
		return nil, apierr.NewInternalServerError()
	}

	var f Bookmark
	for idx, b := range books {
		if b.IsFolder {
			match, err := regexp.MatchString(`(?i)`+folderName, b.Name)
			if err != nil {
				s.log.Errorf("error matching string in regex")
				return nil, apierr.NewInternalServerError()
			}
			if match {
				f = b
				end := len(books) - 1
				books[idx], books[end] = books[end], books[idx]
				books = books[:end]
				break
			}
		}
	}
	folder := organizeBookmarks(books, f.ID, f.Name, f.Path, updatePath(f.Path, f.Name))
	return folder, err
}

// AddBookmark adds a bookmark for an account.
func (s *service) AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate ADD BOOKMARK request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, apierr.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.AddBookmark(reqCtx, requestData, APIKey)
	return numUpdated, err
}

func (s *service) AddBookmarksFromFile(ctx context.Context, r *http.Request, APIKey string) (int, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	header, ok := r.MultipartForm.File[BookmarksFileKey]
	if !ok || len(header) != 1 {
		s.log.Error("Could not find bookmarks_file in request")
		return 0, apierr.NewBadRequestError("no bookmark file in request")
	}
	file, err := header[0].Open()
	if err != nil {
		s.log.Error("Could not open open bookmarks_file")
		return 0, apierr.NewInternalServerError()
	}
	defer file.Close()
	bookmarks, err := NewHTMLBookmarkParser(file, APIKey).parseBookmarkFileHTML()
	if err != nil {
		s.log.Error("Could not parse bookmarks_file")
		return 0, apierr.NewBadRequestError("could not parse bookmark file")
	}
	numAdded, apierr := s.db.AddManyBookmarks(reqCtx, bookmarks)
	if err != nil {
		s.log.Error("Could not add bookmarks to db")
		return 0, apierr
	}
	return numAdded, nil
}

// DeleteBookmark removes a bookmark from an account.
func (s *service) DeleteBookmark(ctx context.Context, bookmarkID, APIKey string) (int, apierr.Error) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Var(bookmarkID, "len=24,hexadecimal")
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate DELETE BOOKMARK request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, apierr.NewBadRequestError("request format incorrect")
	}
	numUpdated, err := s.db.DeleteBookmark(reqCtx, bookmarkID, APIKey)
	return numUpdated, err
}
