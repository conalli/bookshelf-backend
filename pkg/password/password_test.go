package password_test

import (
	"testing"

	"github.com/conalli/bookshelf-backend/pkg/password"
)

func TestHashPassword(t *testing.T) {
	tp := []struct {
		name     string
		password string
	}{
		{
			"simple string",
			"password",
		},
		{
			"long string",
			"qwertyuiopsdfghjklzxcvbnmdfghjklvbnmbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm",
		},
		{
			"mix caps",
			"DFGHJuuioshHHSjeiJOdnuUsizZ",
		},
		{
			"number string",
			"1234567890",
		},
		{
			"non-alphanumeric",
			",./;'[]   -¥",
		},
		{
			"initial space",
			"           password           ",
		},
	}

	for _, p := range tp {
		t.Run(p.name, func(t *testing.T) {
			pw, err := password.HashPassword(p.password)
			if err != nil {
				t.Fatalf("error when attempting to hash password %s", p.password)
			}
			if pw == p.password {
				t.Fatalf("hashed password: %s returned failed to return a hashed result, %s == %s", pw, p.password, pw)
			}
		})
	}
}

func TestCheckHashedPassword(t *testing.T) {
	tc := []struct {
		name     string
		password string
		hashed   string
	}{
		{
			"simple string",
			"password",
			"placeholder",
		},
		{
			"long string",
			"qwertyuiopsdfghjklzxcvbnmdfghjklvbnmbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm",
			"",
		},
		{
			"mix caps",
			"DFGHJuuioshHHSjeiJOdnuUsizZ",
			"",
		},
		{
			"number string",
			"1234567890",
			"",
		},
		{
			"non-alphanumeric",
			",./;'[]   -¥",
			"",
		},
		{
			"initial space",
			"           password           ",
			"",
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			pw, err := password.HashPassword(c.password)
			if err != nil {
				t.Fatalf("error when attempting to hash password %s", c.password)
			}
			c.hashed = pw
			checked := password.CheckHashedPassword(c.hashed, c.password)
			if !checked {
				t.Fatalf("hashed password: %s returned failed to return a hashed result, %s == %s", pw, c.password, pw)
			}
		})
	}
}
