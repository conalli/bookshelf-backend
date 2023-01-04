package accounts

import (
	"github.com/google/uuid"
)

// User represents the db fields associated each user.
type User struct {
	ID            string            `json:"id" bson:"_id,omitempty"`
	APIKey        string            `json:"api_key" bson:"api_key"`
	Name          string            `json:"name" bson:"name"`
	Password      string            `json:"-" bson:"password,omitempty"`
	GivenName     string            `json:"given_name" bson:"given_name"`
	FamilyName    string            `json:"family_name" bson:"family_name"`
	PictureURL    string            `json:"picture" bson:"profile_picture"`
	Email         string            `json:"email" bson:"email"`
	EmailVerified bool              `json:"email_verified" bson:"email_verified"`
	Locale        string            `json:"locale" bson:"locale"`
	Provider      string            `json:"provider" bson:"provider"`
	Cmds          map[string]string `json:"cmds,omitempty" bson:"cmds"`
	Teams         map[string]string `json:"teams,omitempty" bson:"teams"`
}

// GenerateAPIKey generates a random URL-safe string of random length for use as an API key.
func GenerateAPIKey() (string, error) {
	key, err := uuid.NewRandom()
	return key.String(), err
}
