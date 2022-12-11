package request_test

import (
	"net/http"
	"testing"

	"github.com/conalli/bookshelf-backend/pkg/http/request"
	"github.com/google/go-cmp/cmp"
)

func TestFindCookies(t *testing.T) {
	tc := []struct {
		name    string
		cookies []*http.Cookie
	}{
		{
			"two cookie",
			[]*http.Cookie{
				{
					Name:  "test_cookie1",
					Value: "test",
				},
				{
					Name:  "test_cookie2",
					Value: "test",
				},
			},
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			names := []string{}
			for _, cookie := range c.cookies {
				names = append(names, cookie.Name)
			}
			cookies, err := request.FindCookies(c.cookies, names...)
			if err != nil {
				t.Fatal(err)
			}
			if len(cookies) != len(names) {
				t.Error("not enough cookies returned")
			}
			for _, name := range names {
				if _, ok := cookies[name]; !ok {
					t.Errorf("%s not in found cookies", name)
				}
				if cookies[name].Value != "test" {
					t.Error(cmp.Diff(cookies[name].Value, "test"))
				}
			}
		})
	}
}
