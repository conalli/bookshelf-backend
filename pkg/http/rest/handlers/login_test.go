package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/errors"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/go-playground/validator/v10"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db, testutils.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		req        request.LogIn
		statusCode int
		res        accounts.User
	}{
		{
			name: "Default user",
			req: request.LogIn{
				Email:    "default_user@bookshelftest.com",
				Password: "password",
			},
			statusCode: 404,
			res: accounts.User{
				ID:     db.Users["1"].ID,
				APIKey: db.Users["1"].APIKey,
			},
		},
		{
			name: "Incorrect user",
			req: request.LogIn{
				Email:    "incorrect_user@bookshelftest.com",
				Password: "password",
			},
			statusCode: 401,
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body, err := testutils.MakeRequestBody(c.req)
			if err != nil {
				t.Fatalf("Couldn't marshal json body to log in.")
			}
			res, err := http.Post(srv.URL+"/api/auth/login", "application/json", body)
			if err != nil {
				t.Fatalf("Couldn't make post request.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected login request %v to give status code %d: got %d", c.req, c.statusCode, res.StatusCode)
			}
			if res.StatusCode == 404 && res.Request.Header.Get("Referer") != srv.URL+"/api/auth/login" {
				if request.FilterCookies(res.Cookies(), auth.BookshelfAccessToken) != nil {
					t.Errorf("Expected access token cookie to be returned upon log in.")
				}
				if request.FilterCookies(res.Cookies(), auth.BookshelfTokenCode) != nil {
					t.Errorf("Expected code token cookie to be returned upon log in.")
				}
				t.Errorf("Expected redirect upon successful login")
			}
			if res.StatusCode >= 400 && res.StatusCode != 404 {
				var response errors.ResError
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Couldn't decode json body upon sign up.")
				}
			}
			res.Body.Close()
		})
	}
}
