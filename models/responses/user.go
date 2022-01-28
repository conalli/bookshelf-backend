package responses

// SignUpResponse represents the fields returned when a user signs up.
type SignUpResponse struct {
	ID     string `json:"id"`
	APIKey string `json:"apiKey"`
}

// LogInResponse represents the data sent back to the client when the user logs in.
type LogInResponse struct {
	ID     string `json:"id"`
	APIKey string `json:"apiKey"`
}

// AddCmdResponse represents the number of commands that have been updated by the setcmd request.
type AddCmdResponse struct {
	CmdsSet int    `json:"cmdsSet"`
	Cmd     string `json:"cmd"`
	URL     string `json:"url"`
}

// DelCmdResponse represents the fields returned by the delcmd endpoint.
type DelCmdResponse struct {
	NumDeleted int    `json:"numDeleted"`
	Cmd        string `json:"cmd"`
}

// DelUserResponse represents the response from the DelUser endpoint.
type DelUserResponse struct {
	Name       string `json:"name"`
	NumDeleted int    `json:"numDeleted"`
}
