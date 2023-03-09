package handlers_test

import (
	"net/http/httptest"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/go-playground/validator/v10"
)

func TestSearch(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	APIURL := srv.URL + "/api/search/"
	for _, usr := range db.Users {
		for k, v := range usr.Cmds {
			res, err := tu.RequestWithCookie("GET", APIURL+k, tu.WithAPIKey(usr.APIKey))
			if err != nil {
				t.Fatalf("Could not create Search request - %v", err)
			}
			defer res.Body.Close()
			url := res.Request.URL.String()
			if url != v {
				t.Errorf("wanted %s: got %s", v, url)
			}
		}
	}
}
