package search

import (
	"strings"
)

func formatURL(url string) string {
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		return url
	}
	return "http://" + url
}
