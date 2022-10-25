package accounts

import (
	"github.com/google/uuid"
)

// User represents the db fields associated with each user.
type User struct {
	ID       string            `json:"id" bson:"_id,omitempty"`
	Name     string            `json:"name" bson:"name"`
	Password string            `json:"password,omitempty" bson:"password"`
	APIKey   string            `json:"APIKey" bson:"APIKey"`
	Cmds     map[string]string `json:"cmds,omitempty" bson:"cmds"`
	Teams    map[string]string `json:"teams,omitempty" bson:"teams"`
}

// GenerateAPIKey generates a random URL-safe string of random length for use as an API key.
func GenerateAPIKey() (string, error) {
	key, err := uuid.NewRandom()
	return key.String(), err
}
