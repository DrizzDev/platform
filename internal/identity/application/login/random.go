package login

// Random supplies cryptographic randomness for the login state, nonce, PKCE
// verifier, and session identity.
type Random interface {
	Bytes(size int) ([]byte, error)
}
