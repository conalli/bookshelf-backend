package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
)

func TestGetBookmarksFolder(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		folder     string
		APIKey     string
		statusCode int
		res        bookmarks.Folder
	}{
		{
			name:       "Default user, correct request",
			folder:     "News",
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
			res: bookmarks.Folder{
				ID:   "newsfolderid",
				Name: "News",
				Path: bookmarks.BookmarksBasePath,
				Bookmarks: []bookmarks.Bookmark{{
					ID:       "c55fdaace3388c2189875fc5",
					APIKey:   "bd1eb780-0124-11ed-b939-0242ac120002",
					Name:     "bbc",
					Path:     ",News,",
					URL:      "bbc.co.uk",
					IsFolder: false,
				}},
			},
		},
	}
	APIURL := srv.URL + "/api/bookmark/folder?name="
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			res, err := tu.RequestWithCookie("GET", APIURL+c.folder, tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatal("Couldn't create request to get bookmarks folder with cookie.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected get bookmarks folder request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			var response bookmarks.Folder
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatal("Couldn't decode json body upon getting bookmarks folder.")
			}
			if !cmp.Equal(response, c.res) {
				t.Error(cmp.Diff(response, c.res))
			}
			res.Body.Close()
		})
	}
}
