package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/go-playground/validator/v10"
)

func TestSignUp(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db, testutils.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		req        request.SignUp
		statusCode int
	}{
		{
			name: "Correct request",
			req: request.SignUp{
				Email:    "correct_request@bookshelftest.com",
				Password: "password",
			},
			statusCode: 200,
		},
		{
			name: "User Already exists",
			req: request.SignUp{
				Email:    "default_user@bookshelftest.com",
				Password: "password",
			},
			statusCode: 400,
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			body, err := testutils.MakeRequestBody(c.req)
			if err != nil {
				t.Fatalf("Couldn't marshal json body to sign up.")
			}
			res, err := http.Post(srv.URL+"/api/auth/signup", "application/json", body)
			if err != nil {
				t.Fatalf("Couldn't make post request.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected sign up request %v to give status code %d: got %d", c.req, c.statusCode, res.StatusCode)
			}
			if res.StatusCode < 400 {
				var usr accounts.User
				err = json.NewDecoder(res.Body).Decode(&usr)
				if err != nil {
					t.Fatalf("Couldn't decode json body upon sign up.")
				}
				if request.FilterCookies(db.Users["1"].APIKey, res.Cookies()) != nil {
					t.Errorf("Expected jwt cookie to be returned upon log in.")
				}
			}
			res.Body.Close()
		})
	}
}
