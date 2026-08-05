package host

import (
	"context"
	"net/url"

	"github.com/DrizzDev/platform/internal/identity/application/interactive"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/browser"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/system"
	"github.com/DrizzDev/platform/internal/platform/configuration/identity"
)

// access is the login runtime. It runs one browser login bounded by the login
// window, then lets the process-scoped observability shut down.
type access struct {
	foundation
}

func (access access) Run(scope context.Context) (login.Result, error) {
	settings, observer, failure := access.provision(scope)
	if failure != nil {
		return login.Result{}, failure
	}
	current := session{observer: observer}
	defer current.shutdown(scope)
	bounded, cancel := context.WithTimeout(scope, window)
	defer cancel()
	flow, failure := access.assemble(bounded, kit{settings: settings.Identity(), observer: observer, manner: method.Browser})
	if failure != nil {
		return login.Result{}, failure
	}
	access.sweep(bounded, observer)
	result := flow.Run(bounded, login.Input{})
	access.sweep(bounded, observer)
	return result, nil
}

func (access access) assemble(scope context.Context, kit kit) (login.Flow, error) {
	publisher, failure := access.custody(scope, kit)
	if failure != nil {
		return login.Flow{}, failure
	}
	authorizer, failure := access.provider(scope, kit)
	if failure != nil {
		return login.Flow{}, failure
	}
	surface, failure := access.loopback(kit.settings)
	if failure != nil {
		return login.Flow{}, failure
	}
	establishment, failure := interactive.New(interactive.Options{
		Authorization: authorizer,
		Browser:       surface,
		Random:        system.Random{},
	})
	if failure != nil {
		return login.Flow{}, failure
	}
	gate, failure := access.authority(kit)
	if failure != nil {
		return login.Flow{}, failure
	}
	return login.New(login.Options{
		Establishment: establishment,
		Publication:   publisher,
		Authority:     gate,
		Clock:         system.Clock{},
	})
}

func (access access) loopback(settings identity.Settings) (browser.Browser, error) {
	parsed, failure := url.Parse(settings.Redirect())
	if failure != nil {
		return browser.Browser{}, failure
	}
	return browser.New(browser.Options{Opener: browser.System{}, Address: parsed.Host, Path: parsed.Path})
}
