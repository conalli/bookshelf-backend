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

func TestAddBookmark(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	APIKey := db.Users["1"].APIKey
	body, err := testutils.MakeRequestBody(request.AddBookmark{
		Name: "yt",
		Path: ",Google,",
		URL:  "https://www.youtube.com",
	})
	if err != nil {
		t.Fatalf("Couldn't create add cmd request body.")
	}
	res, err := testutils.RequestWithCookie("POST", srv.URL+"/api/user/bookmark/"+APIKey, body, APIKey)
	if err != nil {
		t.Fatalf("Couldn't create request to add cmd with cookie.")
	}
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected add cmd request to give status code %d: got %d", want, res.StatusCode)
	}
	defer res.Body.Close()
	var response handlers.AddBookmarkResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon adding cmds.")
	}
	if response.NumAdded != 1 || response.Name != "yt" || response.URL != "https://www.youtube.com" {
		t.Errorf("Expected commands for user %s to be %v: got %v", db.Users["1"].Name, db.Users["1"].Cmds, response)
	}
}
