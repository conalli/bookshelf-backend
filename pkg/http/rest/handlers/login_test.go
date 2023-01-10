package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/apierr"
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
		err        apierr.ResError
	}{
		{
			name: "Default user",
			req: request.LogIn{
				Email:    "default_user@bookshelftest.com",
				Password: "password",
			},
			statusCode: 200,
			res:        db.Users["1"],
		},
		{
			name: "Incorrect user",
			req: request.LogIn{
				Email:    "incorrect_user@bookshelftest.com",
				Password: "password",
			},
			statusCode: 401,
			err: apierr.ResError{
				Status: 401,
				Title:  apierr.ErrWrongCredentials.Error(),
			},
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
			defer res.Body.Close()
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected login request %v to give status code %d: got %d", c.req, c.statusCode, res.StatusCode)
			}
			if res.StatusCode == 200 {
				if request.FilterCookies(res.Cookies(), auth.BookshelfAccessToken) == nil {
					t.Errorf("Expected access token cookie to be returned upon log in.")
				}
				if request.FilterCookies(res.Cookies(), auth.BookshelfTokenCode) == nil {
					t.Errorf("Expected code token cookie to be returned upon log in.")
				}
				var response accounts.User
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Couldn't decode json body upon getting cmds.")
				}
				if !testutils.IsSameUser(response, c.res) {
					t.Errorf("Expected user to be %v: got %v", c.res, response)
				}
			}
			if res.StatusCode >= 400 {
				var errorResponse apierr.ResError
				err = json.NewDecoder(res.Body).Decode(&errorResponse)
				if err != nil {
					t.Fatalf("Couldn't decode json body upon sign up.")
				}
				if errorResponse.Status != c.err.Status {
					t.Errorf("expected error json to return status code: %d, got: %d", c.err.Status, errorResponse.Status)
				}
				if errorResponse.Title != c.err.Title {
					t.Errorf("expected error json to return error message: %s, got: %s", c.err.Title, errorResponse.Title)
				}
			}
		})
	}
}
