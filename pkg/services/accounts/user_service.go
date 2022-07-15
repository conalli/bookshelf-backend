package accounts

import (
	"context"
	"log"
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/logs"
	"github.com/conalli/bookshelf-backend/pkg/password"
	"github.com/go-playground/validator/v10"
)

// UserRepository provides access to the user storage.
type UserRepository interface {
	NewUser(ctx context.Context, requestData request.SignUp) (User, errors.APIErr)
	GetUserByName(ctx context.Context, requestData request.LogIn) (User, error)
	// GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr)
	AddCmd(reqCtx context.Context, requestData request.AddCmd, APIKey string) (int, errors.APIErr)
	DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, errors.APIErr)
	GetAllBookmarks(ctx context.Context, APIKey string) ([]Bookmark, errors.APIErr)
	GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]Bookmark, errors.APIErr)
	AddBookmark(reqCtx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr)
	DeleteBookmark(reqCtx context.Context, requestData request.DeleteBookmark, APIKey string) (int, errors.APIErr)
	Delete(reqCtx context.Context, requestData request.DeleteUser, APIKey string) (int, errors.APIErr)
}

// UserService provides the user operations.
type UserService interface {
	NewUser(ctx context.Context, requestData request.SignUp) (User, errors.APIErr)
	LogIn(ctx context.Context, requestData request.LogIn) (User, errors.APIErr)
	// GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr)
	GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr)
	AddCmd(reqCtx context.Context, requestData request.AddCmd, APIKey string) (int, errors.APIErr)
	DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, errors.APIErr)
	GetAllBookmarks(ctx context.Context, APIKey string) ([]Bookmark, errors.APIErr)
	GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]Bookmark, errors.APIErr)
	AddBookmark(reqCtx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr)
	DeleteBookmark(reqCtx context.Context, requestData request.DeleteBookmark, APIKey string) (int, errors.APIErr)
	Delete(ctx context.Context, requestData request.DeleteUser, APIKey string) (int, errors.APIErr)
}

type userService struct {
	log      logs.Logger
	validate *validator.Validate
	db       UserRepository
}

// NewUserService creates a search service with the necessary dependencies.
func NewUserService(l logs.Logger, v *validator.Validate, r UserRepository) UserService {
	return &userService{l, v, r}
}

// Search returns the url of a given cmd.
func (s *userService) NewUser(ctx context.Context, requestData request.SignUp) (User, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Struct(requestData)
	if validateErr != nil {
		s.log.Errorf("Could not validate SIGN UP request: %v", validateErr)
		return User{}, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.NewUser(reqCtx, requestData)
	return user, err
}

// Login takes in request data, checks the db and returns the username and apikey is successful.
func (s *userService) LogIn(ctx context.Context, requestData request.LogIn) (User, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Struct(requestData)
	if validateErr != nil {
		s.log.Errorf("Could not validate LOG IN request: %v", validateErr)
		return User{}, errors.NewBadRequestError("request format incorrect.")
	}
	usr, err := s.db.GetUserByName(reqCtx, requestData)
	if err != nil || !password.CheckHashedPassword(usr.Password, requestData.Password) {
		log.Printf("login getuserbykey %+v", err)
		return User{}, errors.NewAPIError(http.StatusUnauthorized, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
	}
	return usr, nil
}

// // GetTeams calls the GetTeams method and returns all teams for a user.
// func (s *userService) GetTeams(ctx context.Context, APIKey string) ([]Team, errors.APIErr) {
// reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
// defer cancelFunc()
// 	teams, err := s.db.GetTeams(reqCtx, APIKey)
// 	return teams, err
// }

// GetAllCmds calls the GetAllCmds method and returns all the users commands.
func (s *userService) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateErr := s.validate.Var(APIKey, "uuid")
	if validateErr != nil {
		s.log.Errorf("Could not validate GET ALL CMDS request: %v", validateErr)
		return nil, errors.NewBadRequestError("request format incorrect.")
	}
	cmds, err := s.db.GetAllCmds(reqCtx, APIKey)
	return cmds, err
}

// AddCmd calls the AddCmd method and returns the number of updated commands.
func (s *userService) AddCmd(ctx context.Context, requestData request.AddCmd, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate ADD CMD request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.AddCmd(reqCtx, requestData, APIKey)
	return numUpdated, err
}

// DeleteCmd calls the DelCmd method and returns the number of updated commands.
func (s *userService) DeleteCmd(ctx context.Context, requestData request.DeleteCmd, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate DELETE CMD request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	numUpdated, err := s.db.DeleteCmd(reqCtx, requestData, APIKey)
	return numUpdated, err
}

func (s *userService) GetAllBookmarks(ctx context.Context, APIKey string) ([]Bookmark, errors.APIErr) {
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

func (s *userService) GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]Bookmark, errors.APIErr) {
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
func (s *userService) AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr) {
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
func (s *userService) DeleteBookmark(ctx context.Context, requestData request.DeleteBookmark, APIKey string) (int, errors.APIErr) {
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

// Delete calls the Delete method and returns the number of deleted users.
func (s *userService) Delete(ctx context.Context, requestData request.DeleteUser, APIKey string) (int, errors.APIErr) {
	reqCtx, cancelFunc := request.CtxWithDefaultTimeout(ctx)
	defer cancelFunc()
	validateReqErr := s.validate.Struct(requestData)
	validateAPIKeyErr := s.validate.Var(APIKey, "uuid")
	if validateReqErr != nil || validateAPIKeyErr != nil {
		s.log.Errorf("Could not validate DELETE USER request: %v - %v", validateReqErr, validateAPIKeyErr)
		return 0, errors.NewBadRequestError("request format incorrect.")
	}
	user, err := s.db.Delete(reqCtx, requestData, APIKey)
	return user, err
}
