package accounts

import "github.com/google/uuid"

// User represents the db fields associated with each user.
type User struct {
	ID       string            `json:"id" bson:"_id,omitempty"`
	Name     string            `json:"name" bson:"name"`
	Password string            `json:"password,omitempty" bson:"password"`
	APIKey   string            `json:"APIKey" bson:"APIKey"`
	Cmds     map[string]string `json:"cmds,omitempty" bson:"cmds"`
	Teams    map[string]string `json:"teams,omitempty" bson:"teams"`
}

// Team represents the db fields associated with each team. Team.Members maps users to roles.
type Team struct {
	ID        string            `json:"id" bson:"_id,omitempty"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	ShortName string            `json:"shortName" bson:"shortName"`
	Cmds      map[string]string `json:"cmds" bson:"cmds"`
	Members   map[string]string `json:"members" bson:"members"`
}

// BookmarkAccount represents an account that owns bookmarks.
type BookmarkAccount struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	APIKey    string     `json:"APIKey" bson:"APIKey"`
	Bookmarks []Bookmark `json:"bookmarks" bson:"bookmarks"`
}

// Bookmark represents a web bookmark.
type Bookmark struct {
	Path string `json:"path" bson:"path,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	URL  string `json:"url" bson:"url"`
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
