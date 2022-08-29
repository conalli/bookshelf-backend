package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/go-playground/validator/v10"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db, testutils.NewCache())
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		body       request.LogIn
		statusCode int
		res        handlers.LogInResponse
	}{
		{
			name: "Default user, correct request",
			body: request.LogIn{
				Name:     "user1",
				Password: "password",
			},
			statusCode: 200,
			res: handlers.LogInResponse{
				ID:     db.Users["1"].ID,
				APIKey: db.Users["1"].APIKey,
			},
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {

			body, err := testutils.MakeRequestBody(c.body)
			if err != nil {
				t.Fatalf("Couldn't marshal json body to log in.")
			}
			res, err := http.Post(srv.URL+"/api/user/login", "application/json", body)
			if err != nil {
				t.Fatalf("Couldn't make post request.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected sign up request %v to give status code %d: got %d", body, c.statusCode, res.StatusCode)
			}
			defer res.Body.Close()
			var response handlers.LogInResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Couldn't decode json body upon sign up.")
			}
			if response.ID != c.res.ID || response.APIKey != c.res.APIKey {
				t.Fatalf("Unexpected log in data")
			}
			if jwtauth.FilterCookies(c.res.APIKey, res.Cookies()) != nil {
				t.Errorf("Expected jwt cookie to be returned upon log in.")
			}
		})
	}
}
