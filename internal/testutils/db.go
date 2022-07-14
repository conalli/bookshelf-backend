package testutils

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/password"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

// Testdb represents a testutils.
type Testdb struct {
	Users map[string]accounts.User
	Teams map[string]accounts.Team
}

// NewDB returns a new Testdb.
func NewDB() *Testdb {
	return &Testdb{}
}

// AddDefaultUsers adds users to an empty testutils.
func (t *Testdb) AddDefaultUsers() *Testdb {
	pw, _ := password.HashPassword("password")
	usrs := map[string]accounts.User{
		"1": {ID: "c55fdaace3388c2189875fc5", Name: "user1", Password: pw, APIKey: "bd1eb780-0124-11ed-b939-0242ac120002", Bookmarks: map[string]string{"bbc": "https://www.bbc.co.uk"}},
	}
	t.Users = usrs
	return t
}

func (t *Testdb) dataAlreadyExists(name string, coll string) bool {
	if coll == "users" {
		for _, v := range t.Users {
			if v.Name == name {
				return true
			}
		}
	}
	if coll == "teams" {
		for _, v := range t.Teams {
			if v.Name == name {
				return true
			}
		}
	}
	return false
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
func (t *Testdb) NewUser(ctx context.Context, body request.SignUp) (accounts.User, errors.APIErr) {
	found := t.dataAlreadyExists(body.Name, "users")
	if found {
		return accounts.User{}, errors.NewBadRequestError("error creating new user; user with name " + body.Name + " already exists")
	}
	key, err := accounts.GenerateAPIKey()
	if err != nil {
		return accounts.User{}, errors.NewInternalServerError()
	}
	usr := accounts.User{
		ID:        body.Name + "999",
		Name:      body.Name,
		Password:  body.Password,
		APIKey:    key,
		Bookmarks: map[string]string{},
		Teams:     map[string]string{},
	}
	t.Users[usr.ID] = usr
	return usr, nil
}

// GetUserByName gets a user by their name in the test db.
func (t *Testdb) GetUserByName(ctx context.Context, body request.LogIn) (accounts.User, error) {
	for _, v := range t.Users {
		if v.Name == body.Name {
			if v.Password == body.Password {
				return accounts.User{}, errors.NewAPIError(403, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
			}
			return v, nil
		}
	}
	return accounts.User{}, errors.NewAPIError(403, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
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
	return usr.Bookmarks, nil
}

// AddCmd adds a cmd to a user in the test db.
func (t *Testdb) AddCmd(reqCtx context.Context, body request.AddCmd, APIKey string) (int, errors.APIErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return 0, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	usr.Bookmarks[body.Cmd] = body.URL
	return 1, nil
}

// DeleteCmd removes a cmd from a user in the test db.
func (t *Testdb) DeleteCmd(ctx context.Context, body request.DeleteCmd, APIKey string) (int, errors.APIErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return 0, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	delete(usr.Bookmarks, body.Cmd)
	return 1, nil
}

// Delete removes a user from the test db.
func (t *Testdb) Delete(reqCtx context.Context, body request.DeleteUser, APIKey string) (int, errors.APIErr) {
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
	val, found := usr.Bookmarks[cmd]
	if !found {
		return "http://www.google.com/search?q=" + cmd, nil
	}
	return val, nil
}
