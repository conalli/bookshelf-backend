package testdb

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

func (t *Testdb) NewUser(ctx context.Context, requestData accounts.SignUpRequest) (accounts.User, errors.ApiErr) {
	found := t.dataAlreadyExists(requestData.Name, "users")
	if found {
		return accounts.User{}, errors.NewBadRequestError("error creating new user; user with name " + requestData.Name + " already exists")
	}
	usr := accounts.User{
		ID:        requestData.Name + "999",
		Name:      requestData.Name,
		Password:  requestData.Password,
		APIKey:    "1234567890",
		Bookmarks: map[string]string{},
		Teams:     map[string]string{},
	}
	t.Users[usr.ID] = usr
	return usr, nil
}

func (t *Testdb) GetUserByName(ctx context.Context, requestData accounts.LogInRequest) (accounts.User, error) {
	for _, v := range t.Users {
		if v.Name == requestData.Name {
			if v.Password == requestData.Password {
				return accounts.User{}, errors.NewApiError(403, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
			}
			return v, nil
		}
	}
	return accounts.User{}, errors.NewApiError(403, errors.ErrWrongCredentials.Error(), "error: name or password incorrect")
}

func (t *Testdb) GetUserByAPIKey(ctx context.Context, APIKey string) (accounts.User, error) {
	for _, v := range t.Users {
		if v.APIKey == APIKey {
			return v, nil
		}
	}
	return accounts.User{}, errors.NewApiError(404, errors.ErrNotFound.Error(), "error: could not find user with APIKey - "+APIKey)
}

// func (t *Testdb) GetTeams(ctx context.Context, APIKey string) ([]accounts.Team, errors.ApiErr) {
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

func (t *Testdb) GetAllCmds(ctx context.Context, APIKey string) (map[string]string, errors.ApiErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return nil, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	return usr.Bookmarks, nil
}

func (t *Testdb) AddCmd(reqCtx context.Context, requestData accounts.AddCmdRequest, APIKey string) (int, errors.ApiErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return 0, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	usr.Bookmarks[requestData.Cmd] = requestData.URL
	return 1, nil
}

func (t *Testdb) DeleteCmd(ctx context.Context, requestData accounts.DelCmdRequest, APIKey string) (int, errors.ApiErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return 0, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	delete(usr.Bookmarks, requestData.Cmd)
	return 1, nil
}

func (t *Testdb) Delete(reqCtx context.Context, requestData accounts.DelUserRequest, APIKey string) (int, errors.ApiErr) {
	usr := t.findUserByAPIKey(APIKey)
	if usr == nil {
		return 0, errors.NewBadRequestError("error: could not find user with value " + APIKey)
	}
	delete(t.Users, requestData.ID)
	return 1, nil
}
