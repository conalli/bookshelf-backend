package auth

const (
	AuthTypeLogIn  string = "login"
	AuthTypeSignUp string = "signup"
)

const (
	ProviderBookshelf string = "bookshelf"
	ProviderGoogle    string = "google"
)

type OIDCRequest struct {
	State, Nonce, AuthURL string
}
