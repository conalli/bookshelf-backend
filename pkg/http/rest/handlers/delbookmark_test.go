package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/go-playground/validator/v10"
)

func TestDeleteBookmark(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db, testutils.NewCache())
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		req        request.DeleteBookmark
		APIKey     string
		statusCode int
		res        handlers.DeleteBookmarkResponse
	}{
		{
			name: "Default bookmark, correct request",
			req: request.DeleteBookmark{
				ID: db.Bookmarks[0].ID,
			},
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
			res: handlers.DeleteBookmarkResponse{
				NumDeleted: 1,
				Name:       db.Bookmarks[0].Name,
				Path:       db.Bookmarks[0].Path,
				URL:        db.Bookmarks[0].URL,
			},
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body, err := testutils.MakeRequestBody(c.req)
			if err != nil {
				t.Fatalf("Couldn't create del bookmark request body.")
			}
			res, err := testutils.RequestWithCookie("DELETE", srv.URL+"/api/bookmark/"+c.APIKey, body, c.APIKey, testutils.NewLogger())
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
			t.Logf("%+v", response)
			if response.NumDeleted != c.res.NumDeleted {
				t.Errorf("Expected %d bookmarks to be deleted: got %d", c.res.NumDeleted, response.NumDeleted)
			}
			res.Body.Close()
		})
	}
}
