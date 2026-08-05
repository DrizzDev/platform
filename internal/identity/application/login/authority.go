package login

import (
	"context"
	"time"
)

// Authority is the login-owned port that authorizes the authenticated subject
// after the credential is published. It resolves the organization with the grant
// the sign-in just obtained — no refresh — and denies a sign-in whose account
// has no usable Drizz access. The cloud remains the authority for every later
// operation.
type Authority interface {
	Authorize(context.Context, Grant) (Tenant, error)
}

// Grant is the freshly issued access token the authority presents to the cloud.
// It is memory-only and never persisted.
type Grant struct {
	Token  []byte
	Expiry time.Time
}

// Tenant is the resolved organization context surfaced to the user on a
// successful sign-in.
type Tenant struct {
	Name string
}
