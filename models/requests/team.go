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

// AddTeamCmdRequest represents the clients request to add a new cmd to a team.
type AddTeamCmdRequest struct {
	ID       string `json:"id"`
	MemberID string `json:"memberId"`
	Cmd      string `json:"cmd"`
	URL      string `json:"url"`
}

// DelTeamCmdRequest represents the expected fields needed for the team/delcmd request to be completed.
type DelTeamCmdRequest struct {
	ID       string `json:"id"`
	MemberID string `json:"memberId"`
	Cmd      string `json:"cmd"`
}
