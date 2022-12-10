package auth

type OIDCRequest struct {
	State, Nonce, AuthURL string
}
