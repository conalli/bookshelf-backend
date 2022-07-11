package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/dbtest"
	"github.com/conalli/bookshelf-backend/internal/handlerstest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest/handlers"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
	"github.com/go-playground/validator/v10"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	db := dbtest.New().AddDefaultUsers()
	r := rest.Router(validator.New(), db, false)
	srv := httptest.NewServer(r)
	defer srv.Close()
	body, err := handlerstest.MakeRequestBody(accounts.LogInRequest{
		Name:     "user1",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("Couldn't marshal json body to log in.")
	}
	res, err := http.Post(srv.URL+"/user/login", "application/json", body)
	if err != nil {
		t.Fatalf("Couldn't make post request.")
	}
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected sign up request %v to give status code %d: got %d", body, want, res.StatusCode)
	}
	defer res.Body.Close()
	var response handlers.LogInResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon sign up.")
	}
	if response.ID != db.Users["1"].ID || response.APIKey != db.Users["1"].APIKey {
		t.Fatalf("Unexpected log in data")
	}
	if jwtauth.FilterCookies(db.Users["1"].APIKey, res.Cookies()) != nil {
		t.Errorf("Expected jwt cookie to be returned upon log in.")
	}
}
