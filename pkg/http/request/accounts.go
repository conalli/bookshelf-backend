package request

// SignUp represents the expected JSON request for the user/signup endpoint.
type SignUp struct {
	Name     string `json:"name" validate:"min=3,max=30"`
	Password string `json:"password" validate:"min=6,max=30"`
}

// LogIn represents the expected JSON request for the user/login endpoint.
type LogIn struct {
	Name     string `json:"name" validate:"min=3,max=30"`
	Password string `json:"password" validate:"min=6,max=30"`
}

// AddCmd represents the expected JSON request for the user/addcmd endpoint.
type AddCmd struct {
	ID  string `json:"id" validate:"len=24,hexadecimal"`
	Cmd string `json:"cmd" validate:"min=1,max=30"`
	URL string `json:"url" validate:"min=5,max=100"`
}

// DeleteCmd represents the expected JSON request for the user/delcmd endpoint.
type DeleteCmd struct {
	ID  string `json:"id" validate:"len=24,hexadecimal"`
	Cmd string `json:"cmd" validate:"min=1,max=30"`
}

// DeleteUser represents the expected JSON request for the user/del endpoint.
type DeleteUser struct {
	ID       string `json:"id" validate:"len=24,hexadecimal"`
	Name     string `json:"name" validate:"min=3,max=30"`
	Password string `json:"password" validate:"min=6,max=30"`
}

// NewTeam represents the expected JSON request for the /team POST endpoint.
type NewTeam struct {
	ID           string `json:"id" validate:"len=24,hexadecimal"`
	Name         string `json:"name" validate:"min=3,max=30"`
	TeamPassword string `json:"password" validate:"min=6,max=30"`
	ShortName    string `json:"shortName"`
}

// DeleteTeam represents the expected JSON request for the /team DELETE endpoint.
type DeleteTeam struct {
	ID           string `json:"id"`
	TeamID       string `json:"teamId"`
	TeamPassword string `json:"password"`
}

// AddMember represents the expected JSON request for the /team/addmember endpoint.
type AddMember struct {
	ID         string `json:"id"`
	TeamID     string `json:"teamId"`
	MemberName string `json:"memberName"`
	Role       string `json:"role"`
}

// DeleteSelf represents the expected JSON request for the /team/delself endpoint.
type DeleteSelf struct {
	ID     string `json:"id"`
	TeamID string `json:"teamId"`
}

// DeleteMember represents the expected JSON request for the /team/delmember endpoint.
type DeleteMember struct {
	ID         string `json:"id"`
	TeamID     string `json:"teamId"`
	MemberName string `jsong:"memberName"`
	Role       string `json:"role"`
}

// AddTeamCmd represents the expected JSON request for the /team/addcmd endpoint.
type AddTeamCmd struct {
	ID       string `json:"id"`
	MemberID string `json:"memberId"`
	Cmd      string `json:"cmd"`
	URL      string `json:"url"`
}

// DeleteTeamCmd represents the expected JSON request for the /team/addcmd endpoint.
type DeleteTeamCmd struct {
	ID       string `json:"id"`
	MemberID string `json:"memberId"`
	Cmd      string `json:"cmd"`
}
