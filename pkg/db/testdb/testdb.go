package testdb

import "github.com/conalli/bookshelf-backend/pkg/accounts"

type testdb struct {
	Users map[string]accounts.User
	Teams map[string]accounts.Team
}

// New returns a new testdb.
func New() *testdb {
	return &testdb{}
}

func (t *testdb) dataAlreadyExists(name string, coll string) bool {
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

func (t *testdb) findUserByAPIKey(APIKey string) *accounts.User {
	for _, v := range t.Users {
		if v.APIKey == APIKey {
			return &v
		}
	}
	return nil
}
