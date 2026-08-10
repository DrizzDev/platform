package cloud

import (
	"context"
	"sync"
	"time"
)

// Provider yields a bearer credential for cloud calls. It is satisfied at the composition root by the identity session
// flow, so capture obtains credentials only through the identity boundary and never mints or stores its own.
type Provider interface {
	Access(context.Context) (Credential, error)
}

// Clock supplies the current time so credential freshness stays deterministic under test.
type Clock interface {
	Now() time.Time
}

// Credential is a short-lived bearer token and the instant it expires.
type Credential struct {
	Token  []byte
	Expiry time.Time
}

func (credential Credential) expired(now time.Time) bool {
	return !now.Before(credential.Expiry)
}

// String redacts the token so an accidental log or error format can never print the secret.
func (Credential) String() string {
	return "cloud.Credential{redacted}"
}

// source caches the credential and renews it through the provider only once it has expired, so draining a batch of
// entries does not re-acquire — and risk rotating — a token on every single upload.
type source struct {
	clock    Clock
	provider Provider
	guard    sync.Mutex
	cached   Credential
}

func (source *source) acquire(scope context.Context) (Credential, error) {
	source.guard.Lock()
	defer source.guard.Unlock()

	if len(source.cached.Token) > 0 && !source.cached.expired(source.clock.Now()) {
		return source.cached, nil
	}
	credential, failure := source.provider.Access(scope)
	if failure != nil {
		return Credential{}, failure
	}
	source.cached = credential
	return credential, nil
}
