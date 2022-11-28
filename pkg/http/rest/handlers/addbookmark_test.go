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
	"github.com/google/uuid"
)

func TestAddBookmark(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db, testutils.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		req        request.AddBookmark
		APIKey     string
		statusCode int
	}{
		{
			name: "Default user",
			req: request.AddBookmark{
				Name: "yt",
				Path: ",Google,",
				URL:  "https://www.youtube.com",
			},
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
		},
		{
			name: "User doesn't exist",
			req: request.AddBookmark{
				Name: "yt",
				Path: ",Google,",
				URL:  "https://www.youtube.com",
			},
			APIKey:     uuid.New().String(),
			statusCode: 400,
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body, err := testutils.MakeRequestBody(c.req)
			if err != nil {
				t.Fatalf("Couldn't create add cmd request body.")
			}
			res, err := testutils.RequestWithCookie("POST", srv.URL+"/api/bookmark/"+c.APIKey, body, c.APIKey, testutils.NewLogger())
			if err != nil {
				t.Fatalf("Couldn't create request to add cmd with cookie.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected add cmd request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			if res.StatusCode < 400 {
				var response handlers.AddBookmarkResponse
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Couldn't decode json body upon adding cmds.")
				}
				if response.NumAdded != 1 {
					t.Errorf("Expected number of commands added for user with API key %s to be 1: got %d", c.APIKey, response.NumAdded)
				}
			}
			res.Body.Close()
		})
	}
}
