package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/go-playground/validator/v10"
)

func TestDeleteBookmark(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		req        string
		APIKey     string
		statusCode int
		res        handlers.DeleteBookmarkResponse
	}{
		{
			name:       "Default bookmark, correct request",
			req:        db.Bookmarks[1].ID,
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
			res: handlers.DeleteBookmarkResponse{
				NumDeleted: 1,
			},
		},
	}
	APIURL := srv.URL + "/api/bookmark/"
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			res, err := tu.RequestWithCookie("DELETE", APIURL+c.req, tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatalf("Couldn't create request to del bookmark with cookie.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected del bookmark request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			var response handlers.DeleteBookmarkResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Couldn't decode json body upon deleting bookmarks.")
			}
			if response.NumDeleted != c.res.NumDeleted {
				t.Errorf("Expected %d bookmarks to be deleted: got %d", c.res.NumDeleted, response.NumDeleted)
			}
			res.Body.Close()
		})
	}
}
