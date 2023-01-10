package bookmarks

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestOrganizeBookmarks(t *testing.T) {
	t.Parallel()
	file, err := os.Open("../../../internal/testdata/bookmarks/safaribookmarks.html")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	APIKey := uuid.New().String()
	bookmarks, err := NewHTMLBookmarkParser(file, APIKey).parseBookmarkFileHTML()
	if err != nil {
		t.Fatal(err)
	}
	got := organizeBookmarks(bookmarks, "id", BookmarksBasePath, BookmarksBasePath, BookmarksBasePath)
	want := &Folder{
		ID:   "id",
		Name: BookmarksBasePath,
		Folders: []Folder{
			{
				Name: "Favourites",
				Bookmarks: []Bookmark{
					{APIKey: APIKey, Name: "Apple", Path: ",Favourites,", URL: "https://www.apple.com/jp/"},
					{APIKey: APIKey, Name: "iCloud", Path: ",Favourites,", URL: "https://www.icloud.com/"},
					{APIKey: APIKey, Name: "Google", Path: ",Favourites,", URL: "https://www.google.co.jp/?client=safari&channel=iphone_bm"},
					{APIKey: APIKey, Name: "Yahoo", Path: ",Favourites,", URL: "https://www.yahoo.co.jp/"},
					{APIKey: APIKey, Name: "Wikipedia", Path: ",Favourites,", URL: "https://ja.wikipedia.org/"},
					{APIKey: APIKey, Name: "Facebook", Path: ",Favourites,", URL: "https://facebook.com/"},
					{APIKey: APIKey, Name: "Twitter", Path: ",Favourites,", URL: "https://twitter.com/"},
					{APIKey: APIKey, Name: "Asahi Shimbun", Path: ",Favourites,", URL: "https://www.asahi.com/"},
					{APIKey: APIKey, Name: "Facebook", Path: ",Favourites,", URL: "https://www.facebook.com/"},
					{APIKey: APIKey, Name: "Wikipedia", Path: ",Favourites,", URL: "http://en.wikipedia.org/wiki/Main_Page"},
					{APIKey: APIKey, Name: "Yahoo!", Path: ",Favourites,", URL: "http://www.yahoo.com/"},
				},
				Folders: []Folder{
					{
						Name: "News",
						Path: ",Favourites,",
						Bookmarks: []Bookmark{
							{APIKey: APIKey, Name: "AllThingsD", Path: ",Favourites,News,", URL: "http://allthingsd.com/"},
							{APIKey: APIKey, Name: "BBC", Path: ",Favourites,News,", URL: "http://www.bbc.co.uk/"},
							{APIKey: APIKey, Name: "CNN", Path: ",Favourites,News,", URL: "http://www.cnn.com/"},
							{APIKey: APIKey, Name: "ESPN", Path: ",Favourites,News,", URL: "http://espn.go.com/"},
							{APIKey: APIKey, Name: "NPR", Path: ",Favourites,News,", URL: "http://www.npr.org/"},
							{APIKey: APIKey, Name: "USA Today", Path: ",Favourites,News,", URL: "http://www.usatoday.com/"},
							{APIKey: APIKey, Name: "The Wall Street Journal", Path: ",Favourites,News,", URL: "http://online.wsj.com/home-page"},
						},
					},
					{
						Name: "Popular",
						Path: ",Favourites,",
						Bookmarks: []Bookmark{
							{APIKey: APIKey, Name: "Amazon", Path: ",Favourites,Popular,", URL: "http://www.amazon.com/"},
							{APIKey: APIKey, Name: "Disney", Path: ",Favourites,Popular,", URL: "http://disney.go.com/"},
							{APIKey: APIKey, Name: "eBay", Path: ",Favourites,Popular,", URL: "http://www.ebay.com/"},
							{APIKey: APIKey, Name: "Flickr", Path: ",Favourites,Popular,", URL: "http://www.flickr.com/"},
							{APIKey: APIKey, Name: "Rotten Tomatoes", Path: ",Favourites,Popular,", URL: "http://www.rottentomatoes.com/"},
							{APIKey: APIKey, Name: "The Weather Channel", Path: ",Favourites,Popular,", URL: "http://www.weather.com/"},
							{APIKey: APIKey, Name: "Yelp", Path: ",Favourites,Popular,", URL: "http://www.yelp.com/"},
							{APIKey: APIKey, Name: "Amazon.co.uk: Low Prices in Electronics, Books, Sports Equipment & more", Path: ",Favourites,Popular,", URL: "http://www.amazon.co.uk/"},
						},
					},
				},
			},
			{Name: "Bookmarks Menu"},
			{Name: "Tab Group Favourites"},
			{Name: "Reading List"},
		},
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
