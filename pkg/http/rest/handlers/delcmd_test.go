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

func TestDeleteCmd(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		req        request.DeleteCmd
		APIKey     string
		statusCode int
		res        handlers.DeleteCmdResponse
	}{
		{
			name: "Default user, correct request.",
			req: request.DeleteCmd{
				ID:  db.Users["1"].ID,
				Cmd: "bbc",
			},
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
			res: handlers.DeleteCmdResponse{
				NumDeleted: 1,
				Cmd:        "bbc",
			},
		},
	}
	APIURL := srv.URL + "/api/user/cmd"
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body, err := tu.MakeJSONRequestBody(c.req)
			if err != nil {
				t.Fatalf("Couldn't create del cmd request body.")
			}
			res, err := tu.RequestWithCookie("PATCH", APIURL, tu.WithBody(body), tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatalf("Couldn't create request to del cmd with cookie.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected del cmd request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			var response handlers.DeleteCmdResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Couldn't decode json body upon deleting cmds.")
			}
			if response.NumDeleted != c.res.NumDeleted || response.Cmd != c.res.Cmd {
				t.Errorf("Expected command %s to be deleted: got %s", c.res.Cmd, response.Cmd)
			}
			res.Body.Close()
		})
	}
}
