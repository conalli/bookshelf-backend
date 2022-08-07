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

func TestDeleteBookmark(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db, testutils.NewCache())
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	APIKey := db.Users["1"].APIKey
	body, err := testutils.MakeRequestBody(request.DeleteBookmark{
		ID: "c55fdaace3388c2189875fc5",
	})
	if err != nil {
		t.Fatalf("Couldn't create del bookmark request body.")
	}
	res, err := testutils.RequestWithCookie("DELETE", srv.URL+"/api/user/bookmark/"+APIKey, body, APIKey, testutils.NewLogger())
	if err != nil {
		t.Fatalf("Couldn't create request to del bookmark with cookie.")
	}
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected del bookmark request to give status code %d: got %d", want, res.StatusCode)
	}
	defer res.Body.Close()
	var response handlers.DeleteBookmarkResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon deleting bookmarks.")
	}
	t.Logf("%+v", response)
	if response.NumDeleted != 1 {
		t.Errorf("Expected %d bookmarks to be deleted: got %d", 1, response.NumDeleted)
	}
}
