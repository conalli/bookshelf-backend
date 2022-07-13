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

func TestDelUser(t *testing.T) {
	t.Parallel()
	db := dbtest.New().AddDefaultUsers()
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	r := rest.NewRouter(sugar, validator.New(), db)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	APIKey := db.Users["1"].APIKey
	body, err := handlerstest.MakeRequestBody(request.DeleteUser{
		ID:       db.Users["1"].ID,
		Name:     db.Users["1"].Name,
		Password: "password",
	})
	if err != nil {
		t.Fatalf("Couldn't create del user request body.")
	}
	res, err := handlerstest.RequestWithCookie("DELETE", srv.URL+"/api/user/"+APIKey, body, APIKey)
	if err != nil {
		t.Fatalf("Couldn't create request to delete user with cookie.")
	}
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected del user request to give status code %d: got %d", want, res.StatusCode)
	}
	defer res.Body.Close()
	var response handlers.DelUserResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon deleting user.")
	}
	if response.NumDeleted != 1 || response.Name != "user1" {
		t.Errorf("Expected NumDeleted to be %d for user %s: got %d", 1, db.Users["1"].Name, response.NumDeleted)
	}
}
