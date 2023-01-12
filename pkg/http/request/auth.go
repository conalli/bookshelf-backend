package request

// SignUp represents the expected JSON request for the user POST endpoint.
type SignUp struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6,max=30"`
}

// LogIn represents the expected JSON request for the user/login POST endpoint.
type LogIn struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6,max=30"`
}
