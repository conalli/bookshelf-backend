package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/go-playground/validator/v10"
)

func TestDeleteCmd(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	APIKey := db.Users["1"].APIKey
	body, err := testutils.MakeRequestBody(request.DeleteCmd{
		ID:  db.Users["1"].ID,
		Cmd: "bbc",
	})
	if err != nil {
		t.Fatalf("Couldn't create del cmd request body.")
	}
	res, err := testutils.RequestWithCookie("PATCH", srv.URL+"/api/user/delcmd/"+APIKey, body, APIKey)
	if err != nil {
		t.Fatalf("Couldn't create request to del cmd with cookie.")
	}
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected del cmd request to give status code %d: got %d", want, res.StatusCode)
	}
	defer res.Body.Close()
	var response handlers.DeleteCmdResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon deleting cmds.")
	}
	if response.NumDeleted != 1 || response.Cmd != "bbc" {
		t.Errorf("Expected commands for user %s to be %v: got %v", db.Users["1"].Name, db.Users["1"].Cmds, response)
	}
}
