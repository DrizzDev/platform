package login

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

// Authorization is the login-owned port over the OAuth/OIDC provider. Begin
// produces the browser redirect; Finish exchanges and validates the callback.
type Authorization interface {
	Begin(context.Context, Secret) (Redirect, error)
	Finish(context.Context, Exchange) (Token, error)
}

// Secret carries the per-attempt values the flow generates and the provider
// binds: CSRF state, OIDC nonce, and the PKCE code verifier.
type Secret struct {
	State    string
	Nonce    string
	Verifier string
}

type Redirect struct {
	URL string
}

type Exchange struct {
	Secret   Secret
	Callback Callback
}

// Token is the validated provider result the flow turns into a credential. The
// Drizz session identity and revision are owned by the flow, not the provider.
type Token struct {
	Issuer  string
	Client  string
	Subject subject.Subject
	Method  method.Method
	Refresh []byte
	Access  []byte
	Issued  time.Time
	Expiry  time.Time
}
