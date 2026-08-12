package host

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/application/reconcile"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	handle "github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/auth0"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/cloud"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/publication"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/system"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/vault"
	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration"
	"github.com/DrizzDev/platform/internal/platform/configuration/identity"
	"github.com/DrizzDev/platform/internal/platform/observability"
	"github.com/DrizzDev/platform/internal/platform/transport"
)

const (
	window     = 5 * time.Minute
	enrollment = 15 * time.Minute
	total      = 15 * time.Second
	connect    = 10 * time.Second
	ceiling    = 1 << 20
	retries    = 2
	minimum    = 200 * time.Millisecond
	maximum    = 2 * time.Second
)

// foundation builds the process-scoped providers and shared identity
// infrastructure that the login and logout runtimes compose on demand.
type foundation struct {
	streams     Streams
	build       build.Info
	environment []string
}

type kit struct {
	settings identity.Settings
	observer observability.Provider
	manner   method.Method
}

func (foundation foundation) provision(scope context.Context) (configuration.Settings, observability.Provider, error) {
	settings, failure := configuration.New(foundation.environment).Load()
	if failure != nil {
		return configuration.Settings{}, observability.Provider{}, failure
	}
	observer, failure := observability.New(scope, observability.Options{
		Build:    foundation.build,
		Settings: settings,
		Output:   foundation.streams.Failure,
	})
	if failure != nil {
		return configuration.Settings{}, observability.Provider{}, failure
	}
	return settings, observer, nil
}

func (foundation foundation) coordination(scope context.Context, observer observability.Provider) (sqlite.Store, error) {
	path, failure := (sqlite.Location{}).Resolve()
	if failure != nil {
		return sqlite.Store{}, failure
	}
	return sqlite.New(scope, sqlite.Options{
		Path:   path,
		Logger: observer.Diagnostics(),
		Tracer: observer.Tracer(),
		Meter:  observer.Meter(),
	})
}

func (foundation foundation) provider(scope context.Context, kit kit) (auth0.Authorizer, error) {
	agent, failure := transport.New(transport.Options{
		Timeout:  total,
		Dial:     connect,
		Ceiling:  ceiling,
		Retries:  retries,
		Minimum:  minimum,
		Maximum:  maximum,
		Tracing:  kit.observer.Tracing(),
		Metering: kit.observer.Metering(),
	})
	if failure != nil {
		return auth0.Authorizer{}, failure
	}
	return auth0.New(scope, auth0.Options{
		Agent:    agent,
		Issuer:   kit.settings.Issuer(),
		Client:   kit.settings.Client(),
		Audience: kit.settings.Audience(),
		Redirect: kit.settings.Redirect(),
		Method:   kit.manner,
		Scopes:   kit.settings.Scopes(),
	})
}

func (foundation foundation) custody(scope context.Context, kit kit) (publication.Publisher, error) {
	store, failure := foundation.coordination(scope, kit.observer)
	if failure != nil {
		return publication.Publisher{}, failure
	}
	safe, failure := vault.New(vault.Options{Store: vault.Keyring{}})
	if failure != nil {
		return publication.Publisher{}, failure
	}
	owner, failure := handle.New(kit.settings.Session())
	if failure != nil {
		return publication.Publisher{}, failure
	}
	return publication.New(publication.Options{Vault: safe, Ledger: store, Random: system.Random{}, Session: owner})
}

// authority builds the post-authentication organization gate. When no cloud
// endpoint is configured it permits every sign-in; otherwise it resolves the
// organization through Drizz Cloud with a fresh, session-backed access token and
// first-party trace propagation.
func (foundation foundation) authority(kit kit) (login.Authority, error) {
	base := kit.settings.Cloud()
	if base == "" {
		return permit{}, nil
	}
	agent, failure := transport.New(transport.Options{
		Timeout:   total,
		Dial:      connect,
		Ceiling:   ceiling,
		Retries:   retries,
		Minimum:   minimum,
		Maximum:   maximum,
		Tracing:   kit.observer.Tracing(),
		Metering:  kit.observer.Metering(),
		Propagate: true,
	})
	if failure != nil {
		return nil, failure
	}
	made, failure := cloud.New(cloud.Options{Agent: agent, Base: base})
	if failure != nil {
		return nil, failure
	}
	return made, nil
}

// sweep drains a bounded batch of the credential cleanup backlog. It runs before
// and after each login and logout so orphaned vault entries are removed promptly
// and a saturated backlog can recover. It is best-effort and never affects the
// operation's outcome.
func (foundation foundation) sweep(scope context.Context, observer observability.Provider) {
	janitor, ready := foundation.janitor(scope, observer)
	if ready {
		janitor.Run(scope, reconcile.Input{})
	}
}

func (foundation foundation) janitor(scope context.Context, observer observability.Provider) (reconcile.Reconciler, bool) {
	store, failure := foundation.coordination(scope, observer)
	if failure != nil {
		return reconcile.Reconciler{}, false
	}
	safe, failure := vault.New(vault.Options{Store: vault.Keyring{}})
	if failure != nil {
		return reconcile.Reconciler{}, false
	}
	made, failure := reconcile.New(reconcile.Options{Queue: store, Vault: safe, Clock: system.Clock{}})
	if failure != nil {
		return reconcile.Reconciler{}, false
	}
	return made, true
}
