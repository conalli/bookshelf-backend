package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/dbtest"
	"github.com/conalli/bookshelf-backend/internal/handlerstest"
	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func TestAddCmd(t *testing.T) {
	t.Parallel()
	db := dbtest.New().AddDefaultUsers()
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	r := rest.NewRouter(sugar, validator.New(), db)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	APIKey := db.Users["1"].APIKey
	body, err := handlerstest.MakeRequestBody(request.AddCmd{
		ID:  db.Users["1"].ID,
		Cmd: "yt",
		URL: "https://www.youtube.com",
	})
	if err != nil {
		t.Fatalf("Couldn't create add cmd request body.")
	}
	res, err := handlerstest.RequestWithCookie("PATCH", srv.URL+"/api/user/addcmd/"+APIKey, body, APIKey)
	if err != nil {
		t.Fatalf("Couldn't create request to add cmd with cookie.")
	}
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected add cmd request to give status code %d: got %d", want, res.StatusCode)
	}
	defer res.Body.Close()
	var response handlers.AddCmdResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon adding cmds.")
	}
	if response.CmdsSet != 1 || response.Cmd != "yt" || response.URL != "https://www.youtube.com" {
		t.Errorf("Expected commands for user %s to be %v: got %v", db.Users["1"].Name, db.Users["1"].Bookmarks, response)
	}
}
