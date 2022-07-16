package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/go-playground/validator/v10"
)

func TestGetBookmarksFolder(t *testing.T) {
	t.Parallel()
	db := testutils.NewDB().AddDefaultUsers()
	r := rest.NewRouter(testutils.NewLogger(), validator.New(), db)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	APIKey := db.Users["1"].APIKey
	res, err := testutils.RequestWithCookie("GET", srv.URL+"/api/user/bookmark/News/"+APIKey, nil, APIKey)
	if err != nil {
		t.Fatalf("Couldn't create request to get bookmarks folder with cookie.")
	}
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected get bookmarks folder request to give status code %d: got %d", want, res.StatusCode)
	}
	defer res.Body.Close()
	var response []accounts.Bookmark
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon getting bookmarks folder.")
	}
	if len(response) != 1 {
		t.Errorf("Expected 1 bookmark for user: got %d", len(response))
	}
}
