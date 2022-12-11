package testutils

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/google/uuid"
)

// Testdb represents a testutils.
type Testdb struct {
	Users     map[string]accounts.User
	Bookmarks []bookmarks.Bookmark
}

// NewDB returns a new Testdb.
func NewDB() *Testdb {
	return &Testdb{}
}

// AddDefaultUsers adds users to an empty testutils.
func (t *Testdb) AddDefaultUsers() *Testdb {
	pw, _ := auth.HashPassword("password")
	usrs := map[string]accounts.User{
		"1": {
			ID:       "c55fdaace3388c2189875fc5",
			Name:     "user1",
			Email:    "default_user@bookshelftest.com",
			Password: pw,
			APIKey:   "bd1eb780-0124-11ed-b939-0242ac120002",
			Cmds:     map[string]string{"bbc": "https://www.bbc.co.uk"},
		},
	}
	t.Users = usrs
	t.Bookmarks = []bookmarks.Bookmark{
		{
			ID:     "c55fdaace3388c2189875fc5",
			APIKey: "bd1eb780-0124-11ed-b939-0242ac120002",
			Name:   "bbc",
			Path:   ",News,",
			URL:    "bbc.co.uk",
		},
	}
	return t
}

func (t *Testdb) UserAlreadyExists(ctx context.Context, email string) (bool, error) {
	for _, v := range t.Users {
		if v.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (t *Testdb) findUserByAPIKey(APIKey string) *accounts.User {
	for _, v := range t.Users {
		if v.APIKey == APIKey {
			return &v
		}
	}
	return nil
}

// NewUser creates a new user in the testdb.
func (t *Testdb) NewUser(ctx context.Context, body request.SignUp) (accounts.User, error) {
	found, _ := t.UserAlreadyExists(context.Background(), body.Email)
	if found {
		return accounts.User{}, errors.ErrBadRequest
	}
	key, err := accounts.GenerateAPIKey()
	if err != nil {
		return accounts.User{}, errors.ErrInternalServerError
	}
	usr := accounts.User{
		ID:       uuid.NewString(),
		Email:    body.Email,
		Password: body.Password,
		APIKey:   key,
		Cmds: map[string]string{
			"yt": "www.youtube.com",
		},
		Teams: map[string]string{},
	}
	t.Users[usr.ID] = usr
	return usr, nil
}

func (t *Testdb) NewOAuthUser(ctx context.Context, IDToken auth.GoogleIDTokenClaims) (accounts.User, error) {
	return accounts.User{}, nil
}

// GetUserByEmail gets a user by their name in the test db.
func (t *Testdb) GetUserByEmail(ctx context.Context, email string) (accounts.User, error) {
	for _, v := range t.Users {
		if v.Email == email {
			return v, nil
		}
	}
	return accounts.User{}, errors.ErrBadRequest
}

// GetUserByAPIKey gets a user by their APIKey in the test db.
func (t *Testdb) GetUserByAPIKey(ctx context.Context, APIKey string) (accounts.User, error) {
	for _, v := range t.Users {
		if v.APIKey == APIKey {
			return v, nil
		}
	}
	return accounts.User{}, errors.NewAPIError(404, errors.ErrNotFound.Error(), "error: could not find user with APIKey - "+APIKey)
}

// func (t *Testdb) GetTeams(ctx context.Context, APIKey string) ([]accounts.Team, errors.APIErr) {
// 	teams := []accounts.Team{}
// 	for _, v := range t.Teams {
// 		for m := range v.Members {
// 			if m == APIKey {
// 				teams = append(teams, v)
// 			}
// 		}
// 	}
// 	return teams, nil
// }

// GetAllCmds gets all cmds for a user in the test db.
func (t *Testdb) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.APIErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return nil, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	return usr.Cmds, nil
}

// AddCmd adds a cmd to a user in the test db.
func (t *Testdb) AddCmd(ctx context.Context, body request.AddCmd, APIKey string) (int, errors.APIErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return 0, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	usr.Cmds[body.Cmd] = body.URL
	return 1, nil
}

// DeleteCmd removes a cmd from a user in the test db.
func (t *Testdb) DeleteCmd(ctx context.Context, body request.DeleteCmd, APIKey string) (int, errors.APIErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return 0, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	delete(usr.Cmds, body.Cmd)
	return 1, nil
}

// GetAllBookmarks gets all bookmarks from the test db.
func (t *Testdb) GetAllBookmarks(ctx context.Context, APIKey string) ([]bookmarks.Bookmark, errors.APIErr) {
	books := make([]bookmarks.Bookmark, 0)
	for _, v := range t.Bookmarks {
		if v.APIKey == APIKey {
			books = append(books, v)
		}
	}
	return books, nil
}

// GetBookmarksFolder gets all bookmarks from the test db.
func (t *Testdb) GetBookmarksFolder(ctx context.Context, path, APIKey string) ([]bookmarks.Bookmark, errors.APIErr) {
	folder := []bookmarks.Bookmark{}
	for _, val := range t.Bookmarks {
		match, err := regexp.Match(path, []byte(val.Path))
		if err != nil {
			return nil, errors.NewBadRequestError("invalid bookmark folder path")
		}
		if match {
			folder = append(folder, val)
		}
	}
	return folder, nil
}

// AddBookmark adds a bookmark to the test db.
func (t *Testdb) AddBookmark(ctx context.Context, requestData request.AddBookmark, APIKey string) (int, errors.APIErr) {
	if _, err := t.GetUserByAPIKey(ctx, APIKey); err != nil {
		return 0, errors.NewBadRequestError("User does not exist.")
	}
	bookmark := bookmarks.Bookmark{
		APIKey: APIKey,
		Name:   requestData.Name,
		Path:   requestData.Path,
		URL:    requestData.URL,
	}
	t.Bookmarks = append(t.Bookmarks, bookmark)
	return 1, nil
}

func (t *Testdb) AddManyBookmarks(ctx context.Context, bookmarks []bookmarks.Bookmark) (int, errors.APIErr) {
	t.Bookmarks = append(t.Bookmarks, bookmarks...)
	return len(bookmarks), nil
}

// DeleteBookmark removes a bookmark from the test db.
func (t *Testdb) DeleteBookmark(ctx context.Context, requestData request.DeleteBookmark, APIKey string) (int, errors.APIErr) {
	i := -1
	for idx := range t.Bookmarks {
		log.Println(idx, t.Bookmarks[idx], t.Bookmarks[idx].ID == requestData.ID)
		if t.Bookmarks[idx].ID == requestData.ID {
			i = idx
			break
		}
	}
	if i < 0 {
		return 0, errors.NewBadRequestError("id not in bookmarks")
	}
	t.Bookmarks[i] = t.Bookmarks[len(t.Bookmarks)-1]
	t.Bookmarks = t.Bookmarks[:len(t.Bookmarks)-1]
	return 1, nil
}

// Delete removes a user from the test db.
func (t *Testdb) Delete(ctx context.Context, body request.DeleteUser, APIKey string) (int, errors.APIErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return 0, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	delete(t.Users, body.ID)
	return 1, nil
}

// Search function for the testutils.
func (t *Testdb) Search(ctx context.Context, APIKey, cmd string) (string, error) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return "", errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	val, found := usr.Cmds[cmd]
	if !found {
		return "http://www.google.com/search?q=" + cmd, nil
	}
	return val, nil
}

// Cache represents a test cache.
type Cache struct {
	Cmds map[string]map[string]string
}

// NewCache returns a new Cache.
func NewCache() *Cache {
	return &Cache{Cmds: map[string]map[string]string{}}
}

// GetSearchData tries to get a URL from the cache.
func (c *Cache) GetSearchData(ctx context.Context, APIKey, cmd string) (string, error) {
	val, ok := c.Cmds[APIKey]
	if !ok {
		return "", fmt.Errorf("no cmds in cache")
	}
	url := val[cmd]
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url, nil
	}
	return "http://" + url, nil
}

// AddCmds adds cmds to the cache.
func (c *Cache) AddCmds(ctx context.Context, APIKey string, cmds map[string]string) bool {
	c.Cmds[APIKey] = cmds
	return true
}

// DeleteCmds removes cmds from the cache.
func (c *Cache) DeleteCmds(ctx context.Context, APIKey string) bool {
	delete(c.Cmds, APIKey)
	return true
}
