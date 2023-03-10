package handlers_test

import (
	"net/http/httptest"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/go-playground/validator/v10"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		APIKey     string
		statusCode int
	}{
		{
			name:       "Correct request, user exists",
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
		},
		{
			name:       "Bad request, user doesn't exist",
			APIKey:     "unknown",
			statusCode: 400,
		},
	}
	APIURL := srv.URL + "/api/user"
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			res, err := tu.RequestWithCookie("GET", APIURL, tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatalf("Couldn't make get request: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected to get status code %d: got %d", c.statusCode, res.StatusCode)
			}
		})
	}
}
