package models

import "strings"

// FormatURL takes a url from the database and returns it in a format required
// for successful redirection.
func FormatURL(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}
	return "http://" + url
}
