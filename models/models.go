package models

type UserData struct {
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	ApiKey    string            `json:"apiKey" bson:"apiKey"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}
