package user

// SignUpRequest represents the expected JSON request for the user/signup endpoint.
type SignUpRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LogInRequest represents the expected JSON request for the user/login endpoint.
type LogInRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// AddCmdRequest represents the expected JSON request for the user/addcmd endpoint.
type AddCmdRequest struct {
	ID  string `json:"id" bson:"_id"`
	Cmd string `json:"cmd" bson:"cmd"`
	URL string `json:"url" bson:"url"`
}

// DelCmdRequest represents the expected JSON request for the user/delcmd endpoint.
type DelCmdRequest struct {
	ID  string `json:"id" bson:"_id"`
	Cmd string `json:"cmd" bson:"cmd"`
}

// DelUserRequest represents the expected JSON request for the user/del endpoint.
type DelUserRequest struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
}
