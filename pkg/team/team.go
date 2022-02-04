package team

// Team represents the db fields associated with each team.
type Team struct {
	ID        string            `json:"id" bson:"_id,omitempty"`
	Name      string            `json:"name" bson:"name"`
	Password  string            `json:"password" bson:"password"`
	ShortName string            `json:"shortName" bson:"shortName"`
	Members   map[string]string `json:"members" bson:"members"`
	Bookmarks map[string]string `json:"bookmarks" bson:"bookmarks"`
}
