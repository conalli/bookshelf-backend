package testdb

import (
	"github.com/conalli/bookshelf-backend/pkg/password"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

// Testdb represents a testdb.
type Testdb struct {
	Users map[string]accounts.User
	Teams map[string]accounts.Team
}

// New returns a new Testdb.
func New() *Testdb {
	return &Testdb{}
}

// AddDefaultUsers adds users to an empty testdb.
func (t *Testdb) AddDefaultUsers() *Testdb {
	pw, _ := password.HashPassword("password")
	usrs := map[string]accounts.User{
		"1": {ID: "1", Name: "user1", Password: pw, APIKey: "111111", Bookmarks: map[string]string{"bbc": "https://www.bbc.co.uk"}},
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
