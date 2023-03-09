package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/go-playground/validator/v10"
)

func TestAddCmd(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		req        request.AddCmd
		APIKey     string
		statusCode int
	}{
		{
			name: "Default User",
			req: request.AddCmd{
				ID:  db.Users["1"].ID,
				Cmd: "yt",
				URL: "https://www.youtube.com",
			},
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
		},
	}
	APIURL := srv.URL + "/api/user/cmd"
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body, err := tu.MakeJSONRequestBody(c.req)
			if err != nil {
				t.Fatalf("Couldn't create add cmd request body.")
			}
			res, err := tu.RequestWithCookie("POST", APIURL, tu.WithBody(body), tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatalf("Couldn't create request to add cmd with cookie.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected add cmd request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			var response handlers.AddCmdResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Couldn't decode json body upon adding cmds.")
			}
			if response.NumAdded != 1 {
				t.Errorf("Expected 1 command for user  with API key %s: got %d", c.APIKey, response.NumAdded)
			}
			res.Body.Close()
		})
	}
}
