package bookmarks

import (
	"os"
	"testing"

	"golang.org/x/net/html"
)

func TestParseBookmarksHTML(t *testing.T) {
	file, err := os.Open("../../../internal/testdata/bookmarks/firefoxbookmarks.html")
	if err != nil {
		t.Fatal(err)
	}
	node, err := html.Parse(file)
	if err != nil {
		t.Fatal(err)
	}
	bookmarks, err := parseBookmarksHTML(node)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bookmarks)
}
