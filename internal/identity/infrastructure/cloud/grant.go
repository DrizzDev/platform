package cloud

import (
	"context"
	"sync"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/grant"
	"github.com/DrizzDev/platform/internal/identity/application/session"
)

// Provider yields a valid access token, renewing the session when the current
// one is stale. It is satisfied by the session flow.
type Provider interface {
	Access(context.Context, session.Input) (grant.Credential, error)
}

// Clock supplies the current time so token freshness stays deterministic.
type Clock interface {
	Now() time.Time
}

// source caches the access token and renews it through the provider only once it
// has expired, so repeated cloud calls do not rotate the refresh token on every
// request.
type source struct {
	provider Provider
	clock    Clock
	guard    sync.Mutex
	cached   grant.Credential
}

func (source *source) acquire(scope context.Context) (grant.Credential, error) {
	source.guard.Lock()
	defer source.guard.Unlock()
	if len(source.cached.Token()) > 0 && !source.cached.Expired(source.clock.Now()) {
		return source.cached, nil
	}
	credential, failure := source.provider.Access(scope, session.Input{})
	if failure != nil {
		return grant.Credential{}, failure
	}
	source.cached = credential
	return credential, nil
}
