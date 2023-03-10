package handlers_test

import (
	"net/http/httptest"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/services/auth"
	"github.com/go-playground/validator/v10"
)

func TestLogOut(t *testing.T) {
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
			name:       "Correct request, user logged in",
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
		},
	}
	APIURL := srv.URL + "/api/auth/logout"
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			res, err := tu.RequestWithCookie("POST", APIURL, tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatalf("Couldn't make post request.")
			}
			if res.StatusCode != c.statusCode {
				t.Errorf("Expected log out request to give status code %d: got %d", c.statusCode, res.StatusCode)
			}
			if res.StatusCode == 200 {
				if bat := request.FilterCookies(res.Cookies(), auth.BookshelfAccessToken); bat.MaxAge != 0 && bat.Value != "" {
					t.Errorf("Expected access token cookie to be deleted upon log out. Cookie: %+v", bat)
				}
				if btc := request.FilterCookies(res.Cookies(), auth.BookshelfTokenCode); btc.MaxAge != 0 && btc.Value != "" {
					t.Errorf("Expected code token cookie to be deleted upon log out. Cookie: %+v", btc)
				}
			}
			res.Body.Close()
		})
	}
}
