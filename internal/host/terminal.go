package host

import (
	"context"

	"github.com/DrizzDev/platform/internal/identity/application/device"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/console"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/system"
)

// terminal is the device-login runtime. It runs one device-authorization grant
// bounded by the enrollment window, presenting the code on the terminal while it
// polls, then publishes the credential like any other sign-in.
type terminal struct {
	foundation
}

func (terminal terminal) Run(scope context.Context) (login.Result, error) {
	settings, observer, failure := terminal.provision(scope)
	if failure != nil {
		return login.Result{}, failure
	}
	current := session{observer: observer}
	defer current.shutdown(scope)
	bounded, cancel := context.WithTimeout(scope, enrollment)
	defer cancel()
	flow, failure := terminal.assemble(bounded, kit{settings: settings.Identity(), observer: observer, manner: method.Device})
	if failure != nil {
		return login.Result{}, failure
	}
	terminal.sweep(bounded, observer)
	result := flow.Run(bounded, login.Input{})
	terminal.sweep(bounded, observer)
	return result, nil
}

func (terminal terminal) assemble(scope context.Context, kit kit) (login.Flow, error) {
	publisher, failure := terminal.custody(scope, kit)
	if failure != nil {
		return login.Flow{}, failure
	}
	authorizer, failure := terminal.provider(scope, kit)
	if failure != nil {
		return login.Flow{}, failure
	}
	screen, failure := console.New(console.Options{Writer: terminal.streams.Output})
	if failure != nil {
		return login.Flow{}, failure
	}
	establishment, failure := device.New(device.Options{Provider: authorizer, Display: screen})
	if failure != nil {
		return login.Flow{}, failure
	}
	gate, failure := terminal.authority(kit)
	if failure != nil {
		return login.Flow{}, failure
	}
	return login.New(login.Options{Establishment: establishment, Publication: publisher, Authority: gate, Clock: system.Clock{}})
}
