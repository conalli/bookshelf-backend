package handlers_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	tu "github.com/conalli/bookshelf-backend/internal/testutils"
	"github.com/conalli/bookshelf-backend/pkg/http/rest"
	"github.com/go-playground/validator/v10"
)

func TestAddBookmarkFile(t *testing.T) {
	t.Parallel()
	db := tu.NewDB().AddDefaultUsers()
	r := rest.NewRouter(tu.NewLogger(), validator.New(), db, tu.NewCache(), nil)
	srv := httptest.NewServer(r.Handler())
	defer srv.Close()
	tc := []struct {
		name       string
		path       string
		APIKey     string
		statusCode int
		want       int
	}{
		{
			name:       "default user, correct request",
			path:       "../../../../internal/testdata/bookmarks/safaribookmarks_basic.html",
			APIKey:     db.Users["1"].APIKey,
			statusCode: 200,
			want:       15,
		},
	}
	APIURL := srv.URL + "/api/bookmark/file"
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			file, ct, err := tu.MakeFileRequestBody(c.path, "safaribookmarks_basic.html")
			if err != nil {
				t.Fatalf("could not create request body: %v", err)
			}
			reqHeaders := map[string]string{
				"Content-Type": ct,
			}
			res, err := tu.RequestWithCookie("POST", APIURL, tu.WithHeaders(reqHeaders), tu.WithBody(file), tu.WithAPIKey(c.APIKey))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if c.statusCode != res.StatusCode {
				t.Errorf("expected status code %d: got %d", c.statusCode, res.StatusCode)
			}
			var got struct {
				NumAdded int `json:"num_added"`
			}
			err = json.NewDecoder(res.Body).Decode(&got)
			if err != nil {
				t.Logf("%+#v", res.Body)
				t.Fatalf("couldn't decode api response: %v", err)
			}
			if c.want != got.NumAdded {
				t.Errorf("wanted: %d, got: %d", c.want, got.NumAdded)
			}
		})
	}
}
