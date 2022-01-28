package requests

// CredentialsRequest represents the fields needed in the request in order to attempt to sign up or log in.
type CredentialsRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// AddCmdRequest represents the expected fields needed for the setcmd request to be completed.
type AddCmdRequest struct {
	ID  string `json:"id" bson:"_id"`
	Cmd string `json:"cmd" bson:"cmd"`
	URL string `json:"url" bson:"url"`
}

// DelCmdRequest represents the expected fields needed for the delcmd request to be completed.
type DelCmdRequest struct {
	ID  string `json:"id" bson:"_id"`
	Cmd string `json:"cmd" bson:"cmd"`
}

// DelUserRequest represents the request body for the DelUser endpoint.
type DelUserRequest struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
}
