package request

// APIRequest represents all API Request types
type APIRequest interface {
	SignUp | LogIn | DeleteUser | AddCmd | DeleteCmd | AddMember
}
