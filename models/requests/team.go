package requests

// NewTeamRequest reprents the clients new team request.
type NewTeamRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	ShortName string `json:"shortName"`
}

// AddMemberRequest represents the clients request to add a new user.
type AddMemberRequest struct {
	ID         string `json:"id"`
	TeamID     string `json:"teamId"`
	MemberName string `json:"memberName"`
	Role       string `json:"role"`
}
