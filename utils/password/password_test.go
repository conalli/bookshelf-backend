package password_test

import (
	"testing"

	"github.com/conalli/bookshelf-backend/utils/password"
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
			"number string",
			"1234567890",
		},
		{
			"long string",
			"qwertyuiopsdfghjklzxcvbnmdfghjklvbnmbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm",
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
