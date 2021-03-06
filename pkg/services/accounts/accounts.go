package accounts

import "github.com/google/uuid"

// User represents the db fields associated with each user.
type User struct {
	ID         string            `json:"id" bson:"_id,omitempty"`
	Name       string            `json:"name" bson:"name"`
	Password   string            `json:"password,omitempty" bson:"password"`
	APIKey     string            `json:"APIKey" bson:"APIKey"`
	BookmarkID string            `json:"bookmarkID" bson:"bookmarkID"`
	Cmds       map[string]string `json:"cmds,omitempty" bson:"cmds"`
	Teams      map[string]string `json:"teams,omitempty" bson:"teams"`
}

// Team represents the db fields associated with each team. Team.Members maps users to roles.
type Team struct {
	ID         string            `json:"id" bson:"_id,omitempty"`
	Name       string            `json:"name" bson:"name"`
	Password   string            `json:"password" bson:"password"`
	ShortName  string            `json:"shortName" bson:"shortName"`
	BookmarkID string            `json:"bookmarkID" bson:"bookmarkID"`
	Cmds       map[string]string `json:"cmds" bson:"cmds"`
	Members    map[string]string `json:"members" bson:"members"`
}

// Bookmark represents a web bookmark.
type Bookmark struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	APIKey string `json:"APIKey" bson:"APIKey"`
	Path   string `json:"path" bson:"path"`
	Name   string `json:"name" bson:"name"`
	URL    string `json:"url" bson:"url"`
}

// Represents the current authorization roles for team members.
const (
	RoleOwner = "OWNER"
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// GenerateAPIKey generates a random URL-safe string of random length for use as an API key.
func GenerateAPIKey() (string, error) {
	key, err := uuid.NewRandom()
	return key.String(), err
}
