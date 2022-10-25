package bookmarks

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"golang.org/x/net/html"
)

func TestParseBookmarksHTML(t *testing.T) {
	file, err := os.Open("../../../internal/testdata/bookmarks/safaribookmarks.html")
	if err != nil {
		t.Fatal(err)
	}
	node, err := html.Parse(file)
	if err != nil {
		t.Fatal(err)
	}
	APIKey := uuid.New().String()
	got, err := parseBookmarksHTML(APIKey, node)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(got)
	want := []Bookmark{
		{APIKey: APIKey, Name: "Apple", Path: ",Favorites,", URL: "https://www.apple.com/jp/"},
		{APIKey: APIKey, Name: "iCloud", Path: ",Favorites,", URL: "https://www.icloud.com/"},
		{APIKey: APIKey, Name: "Google", Path: ",Favorites,", URL: "https://www.google.co.jp/?client=safari&channel=iphone_bm"},
		{APIKey: APIKey, Name: "Yahoo", Path: ",Favorites,", URL: "https://www.yahoo.co.jp/"},
		{APIKey: APIKey, Name: "Wikipedia", Path: ",Favorites,", URL: "https://ja.wikipedia.org/"},
		{APIKey: APIKey, Name: "Facebook", Path: ",Favorites,", URL: "https://facebook.com/"},
		{APIKey: APIKey, Name: "Twitter", Path: ",Favorites,", URL: "https://twitter.com/"},
		{APIKey: APIKey, Name: "Asahi Shimbun", Path: ",Favorites,", URL: "https://www.asahi.com/"},
		{APIKey: APIKey, Name: "Facebook", Path: ",Favorites,", URL: "https://www.facebook.com/"},
		{APIKey: APIKey, Name: "Wikipedia", Path: ",Favorites,", URL: "http://en.wikipedia.org/wiki/Main_Page"},
		{APIKey: APIKey, Name: "Yahoo!", Path: ",Favorites,", URL: "http://www.yahoo.com/"},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
