package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/conalli/bookshelf-backend/internal/dbtest"
	"github.com/conalli/bookshelf-backend/internal/handlerstest"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/conalli/bookshelf-backend/pkg/jwtauth"
	"github.com/conalli/bookshelf-backend/pkg/services/accounts"
)

func TestSignUp(t *testing.T) {
	t.Parallel()
	db := dbtest.New().AddDefaultUsers()
	r := rest.Router(db, false)
	srv := httptest.NewServer(r)
	defer srv.Close()
	body, err := handlerstest.MakeRequestBody(accounts.SignUpRequest{
		Name:     "signuptest",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("Couldn't marshal json body to sign up.")
	}
	res, err := http.Post(srv.URL+"/user", "application/json", body)
	if err != nil {
		t.Fatalf("Couldn't make post request.")
	}
	defer res.Body.Close()
	want := 200
	if res.StatusCode != want {
		t.Errorf("Expected sign up request %v to give status code %d: got %d", body, want, res.StatusCode)
	}
	var usr accounts.User
	err = json.NewDecoder(res.Body).Decode(&usr)
	if err != nil {
		t.Fatalf("Couldn't decode json body upon sign up.")
	}
	// Change password hashing logic
	if usr.ID != usr.Name+"999" || usr.Name != "signuptest" || usr.APIKey != "1234567890" || usr.Password != "password" {
		t.Fatalf("Unexpected sign up data")
	}
	if jwtauth.FilterCookies(db.Users["1"].APIKey, res.Cookies()) != nil {
		t.Errorf("Expected jwt cookie to be returned upon log in.")
	}
	// Check user already exists case
}
