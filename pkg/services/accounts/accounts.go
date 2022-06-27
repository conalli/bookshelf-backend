package accounts

import "github.com/google/uuid"

// User represents the db fields associated with each user.
type User struct {
	ID        string            `json:"id" bson:"_id,omitempty"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password,omitempty" bson:"password"`
	APIKey    string            `json:"APIKey" bson:"APIKey"`
	Bookmarks map[string]string `json:"bookmarks,omitempty" bson:"bookmarks"`
	Teams     map[string]string `json:"teams,omitempty" bson:"teams"`
}

// Team represents the db fields associated with each team. Team.Members maps users to roles.
type Team struct {
	ID        string            `json:"id" bson:"_id,omitempty"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	ShortName string            `json:"shortName" bson:"shortName"`
	Members   map[string]string `json:"members" bson:"members"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}

// Represents the current authorization roles for team members.
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// GenerateAPIKey generates a random URL-safe string of random length for use as an API key.
func GenerateAPIKey() (string, error) {
	key, err := uuid.NewRandom()
	return key.String(), err
}
