package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/apierr"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func TestAddBookmark(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
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
				Name:     "yt",
				Path:     ",Google,",
				URL:      "https://www.youtube.com",
				IsFolder: false,
			},
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
		},
		{
			name: "User doesn't exist",
			req: request.AddBookmark{
				Name:     "yt",
				Path:     ",Google,",
				URL:      "https://www.youtube.com",
				IsFolder: false,
			},
			APIKey:     uuid.New().String(),
			statusCode: 400,
		},
	}
	APIURL := srv.URL + "/api/bookmark"
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body, err := tu.MakeJSONRequestBody(c.req)
			if err != nil {
				t.Fatalf("Couldn't create add cmd request body")
			}
			res, err := tu.RequestWithCookie("POST", APIURL, tu.WithBody(body), tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatalf("Couldn't create request to add cmd with cookie")
			}
			defer res.Body.Close()
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected add cmd request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			if res.StatusCode < 400 {
				var response handlers.AddBookmarkResponse
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Couldn't decode json body upon adding cmds")
				}
				if response.NumAdded != 1 {
					t.Errorf("Expected number of commands added for user with API key %s to be 1: got %d", c.APIKey, response.NumAdded)
				}
			} else {
				var response apierr.ResError
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Couldn't decode json body upon add cmd error")
				}
			}
		})
	}
}
