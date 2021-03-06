package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/go-playground/validator/v10"
)

func TestGetCmds(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	APIKey := db.Users["1"].APIKey
	res, err := testutils.RequestWithCookie("GET", srv.URL+"/api/user/cmd/"+APIKey, nil, APIKey)
	if err != nil {
		t.Fatalf("Couldn't create request to get cmds with cookie.")
	}
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected get cmd request to give status code %d: got %d", want, res.StatusCode)
	}
	defer res.Body.Close()
	var response map[string]string
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon getting cmds.")
	}
	if fmt.Sprint(response) != fmt.Sprint(db.Users["1"].Cmds) {
		t.Errorf("Expected commands for user %s to be %v: got %v", db.Users["1"].Name, db.Users["1"].Cmds, response)
	}
}
