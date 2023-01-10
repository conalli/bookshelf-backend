package accounts

// User represents the db fields associated each user.
type User struct {
	ID            string            `json:"id" bson:"_id,omitempty" redis:"id"`
	APIKey        string            `json:"api_key" bson:"api_key" redis:"api_key"`
	Name          string            `json:"name" bson:"name" redis:"name"`
	Password      string            `json:"-" bson:"password,omitempty"`
	GivenName     string            `json:"given_name" bson:"given_name" redis:"given_name"`
	FamilyName    string            `json:"family_name" bson:"family_name" redis:"family_name"`
	PictureURL    string            `json:"picture" bson:"profile_picture" redis:"picture"`
	Email         string            `json:"email" bson:"email" redis:"email"`
	EmailVerified bool              `json:"email_verified" bson:"email_verified" redis:"email_verified"`
	Locale        string            `json:"locale" bson:"locale" redis:"locale"`
	Provider      string            `json:"provider" bson:"provider" redis:"provider"`
	Cmds          map[string]string `json:"cmds,omitempty" bson:"cmds"`
	Teams         map[string]string `json:"teams,omitempty" bson:"teams"`
}
