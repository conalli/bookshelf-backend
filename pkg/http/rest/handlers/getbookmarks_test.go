package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/services/bookmarks"
	"github.com/go-playground/validator/v10"
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
		res        []bookmarks.Bookmark
	}{
		{
			name:       "Default user, correct request",
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
			res:        []bookmarks.Bookmark{db.Bookmarks[0]},
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
			var response []bookmarks.Bookmark
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Couldn't decode json body upon getting bookmarks.")
			}
			if len(response) != len(c.res) {
				t.Errorf("Expected 1 bookmark for user: got %d", len(response))
			}
			res.Body.Close()
		})
	}
}
