package requests

// NewTeamReq reprents the clients new team request.
type NewTeamReq struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	ShortName string `json:"shortName"`
}

// AddMemberReq represents the clients request to add a new user.
type AddMemberReq struct {
	ID         string `json:"id"`
	TeamID     string `json:"teamId"`
	MemberName string `json:"memberName"`
	Role       string `json:"role"`
}
