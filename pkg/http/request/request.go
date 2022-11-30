package request

import "net/http"

// APIRequest represents all API Request types
type APIRequest interface {
	SignUp | LogIn | DeleteUser | AddCmd | DeleteCmd | AddBookmark | DeleteBookmark
}

// FilterCookies looks through all cookies and returns one with given name.
func FilterCookies(name string, cookies []*http.Cookie) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}
