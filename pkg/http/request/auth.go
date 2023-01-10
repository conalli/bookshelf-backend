package request

import "net/http"

const JWTAPIKey ContextKey = "api_key"

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

func GetAPIKeyFromContext(r *http.Request) (string, bool) {
	key, ok := r.Context().Value(JWTAPIKey).(string)
	return key, ok
}
