package host

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/logout"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/system"
)

// departure is the logout runtime. It clears the local session and attempts
// bounded server revocation, then lets observability shut down.
type departure struct {
	foundation
}

func (departure departure) Run(scope context.Context) (logout.Result, error) {
	settings, observer, failure := departure.provision(scope)
	if failure != nil {
		return logout.Result{}, failure
	}
	current := session{observer: observer}
	defer current.shutdown(scope)
	bounded, cancel := context.WithTimeout(scope, window)
	defer cancel()
	flow, failure := departure.assemble(bounded, kit{settings: settings.Identity(), observer: observer, manner: method.Browser})
	if failure != nil {
		return logout.Result{}, failure
	}
	departure.sweep(bounded, observer)
	result := flow.Run(bounded, logout.Input{})
	departure.sweep(bounded, observer)
	return result, nil
}

func (departure departure) assemble(scope context.Context, kit kit) (logout.Flow, error) {
	publisher, failure := departure.custody(scope, kit)
	if failure != nil {
		return logout.Flow{}, failure
	}
	authorizer, failure := departure.provider(scope, kit)
	if failure != nil {
		return logout.Flow{}, failure
	}
	return logout.New(logout.Options{
		Vault:       publisher,
		Publication: publisher,
		Revocation:  authorizer,
		Clock:       system.Clock{},
	})
}
