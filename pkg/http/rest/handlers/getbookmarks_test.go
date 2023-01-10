package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
)

func TestGetBookmarks(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db, testutils.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		APIKey     string
		statusCode int
		res        bookmarks.Folder
	}{
		{
			name:       "Default user, correct request",
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
			res: bookmarks.Folder{
				ID:        "",
				Name:      "",
				Path:      "",
				Bookmarks: nil,
				Folders: []bookmarks.Folder{
					{
						ID:        "newsfolderid",
						Name:      "News",
						Bookmarks: []bookmarks.Bookmark{db.Bookmarks[1]},
					},
				},
			},
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			res, err := testutils.RequestWithCookie("GET", srv.URL+"/api/bookmark", nil, c.APIKey, testutils.NewLogger())
			if err != nil {
				t.Fatalf("Couldn't create request to get bookmarks with cookie.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected get bookmarks request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			var response bookmarks.Folder
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Couldn't decode json body upon getting bookmarks.")
			}
			if !cmp.Equal(response, c.res) {
				t.Errorf(cmp.Diff(response, c.res))
			}
			res.Body.Close()
		})
	}
}
