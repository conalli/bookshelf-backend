package accounts

// SignUpRequest represents the expected JSON request for the user/signup endpoint.
type SignUpRequest struct {
	Name     string `json:"name" validate:"min=3,max=30"`
	Password string `json:"password" validate:"min=6,max=30"`
}

// LogInRequest represents the expected JSON request for the user/login endpoint.
type LogInRequest struct {
	Name     string `json:"name" validate:"min=3,max=30"`
	Password string `json:"password" validate:"min=6,max=30"`
}

// AddCmdRequest represents the expected JSON request for the user/addcmd endpoint.
type AddCmdRequest struct {
	ID  string `json:"id" validate:"len=24,hexadecimal"`
	Cmd string `json:"cmd" validate:"min=1,max=30"`
	URL string `json:"url" validate:"min=5,max=100"`
}

// DelCmdRequest represents the expected JSON request for the user/delcmd endpoint.
type DelCmdRequest struct {
	ID  string `json:"id" validate:"len=24,hexadecimal"`
	Cmd string `json:"cmd" validate:"min=1,max=30"`
}

// DelUserRequest represents the expected JSON request for the user/del endpoint.
type DelUserRequest struct {
	ID       string `json:"id" validate:"len=24,hexadecimal"`
	Name     string `json:"name" validate:"min=3,max=30"`
	Password string `json:"password" validate:"min=6,max=30"`
}

// NewTeamRequest represents the expected JSON request for the /team POST endpoint.
type NewTeamRequest struct {
	ID           string `json:"id" validate:"len=24,hexadecimal"`
	Name         string `json:"name" validate:"min=3,max=30"`
	TeamPassword string `json:"password" validate:"min=6,max=30"`
	ShortName    string `json:"shortName"`
}

// DelTeamRequest represents the expected JSON request for the /team DELETE endpoint.
type DelTeamRequest struct {
	ID           string `json:"id"`
	TeamID       string `json:"teamId"`
	TeamPassword string `json:"password"`
}

// AddMemberRequest represents the expected JSON request for the /team/addmember endpoint.
type AddMemberRequest struct {
	ID         string `json:"id"`
	TeamID     string `json:"teamId"`
	MemberName string `json:"memberName"`
	Role       string `json:"role"`
}

// DelSelfRequest represents the expected JSON request for the /team/delself endpoint.
type DelSelfRequest struct {
	ID     string `json:"id"`
	TeamID string `json:"teamId"`
}

// DelMemberRequest represents the expected JSON request for the /team/delmember endpoint.
type DelMemberRequest struct {
	ID         string `json:"id"`
	TeamID     string `json:"teamId"`
	MemberName string `jsong:"memberName"`
	Role       string `json:"role"`
}

// AddTeamCmdRequest represents the expected JSON request for the /team/addcmd endpoint.
type AddTeamCmdRequest struct {
	ID       string `json:"id"`
	MemberID string `json:"memberId"`
	Cmd      string `json:"cmd"`
	URL      string `json:"url"`
}

// DelTeamCmdRequest represents the expected JSON request for the /team/addcmd endpoint.
type DelTeamCmdRequest struct {
	ID       string `json:"id"`
	MemberID string `json:"memberId"`
	Cmd      string `json:"cmd"`
}
