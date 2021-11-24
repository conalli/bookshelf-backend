package password

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tp := []struct {
		name     string
		password string
		expect   bool
	}{
		{
			"simple string",
			"password",
			true,
		},
		{
			"number string",
			"1234567890",
			true,
		},
		{
			"long string",
			"qwertyuiopsdfghjklzxcvbnmdfghjklvbnmbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm",
			true,
		},
	}

	for _, p := range tp {
		t.Run(p.name, func(t *testing.T) {
			pw, err := HashPassword(p.password)
			if err != nil {
				t.Fatalf("error when attempting to hash password %s", p.password)
			}
			if pw == p.password {
				t.Fatalf("hashed password: %s returned failed to return a hashed result, %s == %s", pw, p.password, pw)
			}
		})
	}
}
