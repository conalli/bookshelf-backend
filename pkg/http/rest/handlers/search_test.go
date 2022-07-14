package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/go-playground/validator/v10"
)

func TestSearch(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	for _, usr := range db.Users {
		for k, v := range usr.Cmds {
			res, err := http.Get(fmt.Sprintf("%s/api/search/%s/%s", srv.URL, usr.APIKey, k))
			if err != nil {
				t.Fatalf("Could not create Search request - %v", err)
			}
			defer res.Body.Close()
			url := res.Request.URL.String()
			if url != v {
				t.Errorf("wanted %s: got %s", v, url)
			}
			res, err = http.Get(fmt.Sprintf("%s/api/search/%s/%s", srv.URL, usr.APIKey, k+"test"))
			if err != nil {
				t.Fatalf("Could not create Search request - %v", err)
			}
			url = res.Request.URL.String()
			google := "http://www.google.com/search?q=" + k + "test"
			if url != google {
				t.Errorf("wanted %s: got %s", v, url)
			}
		}
	}
}
