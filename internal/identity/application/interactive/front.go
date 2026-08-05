package interactive

import (
	"context"
	"encoding/base64"

	"github.com/DrizzDev/platform/internal/identity/application/login"
)

const entropy = 32

var _ login.Establishment = front{}

// front establishes a sign-in through the system browser using Authorization
// Code with PKCE. It generates the per-attempt CSRF state, OIDC nonce, and PKCE
// verifier, opens the redirect, and exchanges the captured callback.
type front struct {
	authorization login.Authorization
	browser       login.Browser
	random        login.Random
}

func (front front) Establish(scope context.Context) (login.Token, error) {
	secret, failure := front.secret()
	if failure != nil {
		return login.Token{}, failure
	}
	redirect, failure := front.authorization.Begin(scope, secret)
	if failure != nil {
		return login.Token{}, failure
	}
	callback, failure := front.browser.Prompt(scope, redirect)
	if failure != nil {
		return login.Token{}, failure
	}
	return front.authorization.Finish(scope, login.Exchange{Secret: secret, Callback: callback})
}

func (front front) secret() (login.Secret, error) {
	state, failure := front.mint()
	if failure != nil {
		return login.Secret{}, failure
	}
	nonce, failure := front.mint()
	if failure != nil {
		return login.Secret{}, failure
	}
	verifier, failure := front.mint()
	if failure != nil {
		return login.Secret{}, failure
	}
	return login.Secret{State: state, Nonce: nonce, Verifier: verifier}, nil
}

func (front front) mint() (string, error) {
	bytes, failure := front.random.Bytes(entropy)
	if failure != nil {
		return "", failure
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
