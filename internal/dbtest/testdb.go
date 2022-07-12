package dbtest

import (
	"github.com/conalli/bookshelf-backend/pkg/password"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

// Testdb represents a dbtest.
type Testdb struct {
	Users map[string]accounts.User
	Teams map[string]accounts.Team
}

// New returns a new Testdb.
func New() *Testdb {
	return &Testdb{}
}

// AddDefaultUsers adds users to an empty dbtest.
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
