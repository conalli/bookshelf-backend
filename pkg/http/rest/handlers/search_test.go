package handlers_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/go-playground/validator/v10"
)

func TestSearchCommand(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	client := tu.NewRedirectClient()
	APIURL := srv.URL + "/api/search/"
	for _, usr := range db.Users {
		for cmd, URL := range usr.Cmds {
			res, err := tu.RequestWithCookie("GET", APIURL+cmd, tu.WithClient(client), tu.WithAPIKey(usr.APIKey))
			if err != nil {
				t.Fatalf("Could not create Search request - %v", err)
			}
			defer res.Body.Close()
			dest := res.Header.Get("Location")
			if dest != URL {
				t.Errorf("wanted %s: got %s", URL, dest)
			}
		}
	}
}

func TestSearchLS(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	redirectURL := os.Getenv("ALLOWED_URL_BASE")
	tc := []struct {
		name        string
		APIKey      string
		flags       string
		statusCode  int
		redirectURL string
	}{
		{
			name:        "Correct request, (ls -b)",
			APIKey:      db.Users["1"].APIKey,
			flags:       "-b",
			statusCode:  303,
			redirectURL: redirectURL + "/webcli/bookmark",
		},
		{
			name:        "Correct request, (ls -c)",
			APIKey:      db.Users["1"].APIKey,
			flags:       "-c",
			statusCode:  303,
			redirectURL: redirectURL + "/webcli/command",
		},
		{
			name:        "Correct request, (ls -bf)",
			APIKey:      db.Users["1"].APIKey,
			flags:       "-bf news",
			statusCode:  303,
			redirectURL: redirectURL + "/webcli/bookmark?folder=news",
		},
		{
			name:        "Incorrect request, incorrect APIKey (ls -b)",
			APIKey:      "unknown",
			flags:       "-b",
			statusCode:  303,
			redirectURL: redirectURL + "/webcli/error",
		},
	}
	APIURL := srv.URL + "/api/search/ls"
	client := tu.NewRedirectClient()
	for _, c := range tc {
		res, err := tu.RequestWithCookie("GET", fmt.Sprintf("%s %s", APIURL, c.flags), tu.WithClient(client), tu.WithAPIKey(c.APIKey))
		if err != nil {
			t.Fatalf("Could not create Search request - %v", err)
		}
		defer res.Body.Close()
		if res.StatusCode != c.statusCode {
			t.Errorf("wanted %d: got %d", c.statusCode, res.StatusCode)
		}
		url := res.Header.Get("Location")
		if url != c.redirectURL {
			t.Errorf("wanted %s: got %s", c.redirectURL, url)
		}
	}
}
