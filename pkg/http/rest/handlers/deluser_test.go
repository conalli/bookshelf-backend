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

func TestDelUser(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		req        request.DeleteUser
		APIKey     string
		statusCode int
		res        handlers.DelUserResponse
	}{
		{
			name: "Default user, correct request.",
			req: request.DeleteUser{
				ID:       db.Users["1"].ID,
				Name:     db.Users["1"].Name,
				Password: "password",
			},
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
			res: handlers.DelUserResponse{
				NumDeleted: 1,
				Name:       db.Users["1"].Name,
			},
		},
	}
	APIURL := srv.URL + "/api/user"
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body, err := tu.MakeJSONRequestBody(c.req)
			if err != nil {
				t.Fatalf("Couldn't create del user request body.")
			}
			res, err := tu.RequestWithCookie("DELETE", APIURL, tu.WithBody(body), tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatalf("Couldn't create request to delete user with cookie.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected del user request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			var response handlers.DelUserResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Couldn't decode json body upon deleting user.")
			}
			if response.NumDeleted != c.res.NumDeleted || response.Name != c.res.Name {
				t.Errorf("Expected NumDeleted to be %d for user %s: got %d", c.res.NumDeleted, c.res.Name, response.NumDeleted)
			}
			res.Body.Close()
		})
	}
}
