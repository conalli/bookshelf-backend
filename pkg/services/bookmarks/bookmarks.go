package bookmarks

import "golang.org/x/net/html"

// Bookmark represents a web bookmark.
type Bookmark struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	APIKey string `json:"APIKey" bson:"APIKey"`
	Path   string `json:"path" bson:"path"`
	Name   string `json:"name" bson:"name"`
	URL    string `json:"url" bson:"url"`
}

func parseBookmarksHTML(root *html.Node) ([]Bookmark, error) {
	return nil, nil
}
