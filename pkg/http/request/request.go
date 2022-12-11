package request

import (
	"net/http"

	"github.com/conalli/bookshelf-backend/pkg/errors"
)

// APIRequest represents all API Request types
type APIRequest interface {
	SignUp | LogIn | DeleteUser | AddCmd | DeleteCmd | AddBookmark | DeleteBookmark
}

// FilterCookies looks through all cookies and returns cookie with given name.
func FilterCookies(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func FindCookies(cookies []*http.Cookie, names ...string) (map[string]*http.Cookie, error) {
	found := map[string]*http.Cookie{}
	for _, name := range names {
		cookie := FilterCookies(cookies, name)
		if cookie == nil {
			return nil, errors.ErrBadRequest
		}
		found[name] = cookie
	}
	return found, nil
}
