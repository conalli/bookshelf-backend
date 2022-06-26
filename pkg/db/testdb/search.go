package testdb

import (
	"context"

	"github.com/conalli/bookshelf-backend/pkg/errors"
)

func (t *testdb) Search(ctx context.Context, APIKey, cmd string) (string, error) {
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
