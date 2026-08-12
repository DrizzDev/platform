package host

import (
	"context"

	integration "github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/binary"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/telemetry"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/wiring"
	"github.com/DrizzDev/platform/internal/transport/cli"
)

var _ cli.Connector = connector{}

// connector is the host's installer runtime. It assembles the installer — the executable resolver, the per-agent
// configuration store, and the recorder that writes each change into the same durable store device executions use —
// on each call, since connecting an agent is a rare, one-shot action rather than a hot path.
type connector struct {
	base foundation
}

func (connector connector) Survey(scope context.Context) integration.Report {
	installer, release, failure := connector.assemble(scope)
	if failure != nil {
		return integration.Report{}
	}
	defer release()
	return installer.Survey(scope)
}

func (connector connector) Enable(scope context.Context, selection integration.Selection) (integration.Report, error) {
	installer, release, failure := connector.assemble(scope)
	if failure != nil {
		return integration.Report{}, failure
	}
	defer release()
	return installer.Enable(scope, selection)
}

func (connector connector) Disable(scope context.Context, selection integration.Selection) integration.Report {
	installer, release, failure := connector.assemble(scope)
	if failure != nil {
		return integration.Report{}
	}
	defer release()
	return installer.Disable(scope, selection)
}

func (connector connector) Capture(scope context.Context, selection integration.Selection) (integration.Report, error) {
	installer, release, failure := connector.assemble(scope)
	if failure != nil {
		return integration.Report{}, failure
	}
	defer release()
	return installer.Capture(scope, selection)
}

func (connector connector) Uncapture(scope context.Context, selection integration.Selection) integration.Report {
	installer, release, failure := connector.assemble(scope)
	if failure != nil {
		return integration.Report{}
	}
	defer release()
	return installer.Uncapture(scope, selection)
}

func (connector connector) assemble(scope context.Context) (integration.Installer, func(), error) {
	var closers []func()
	release := func() {
		for index := len(closers) - 1; index >= 0; index-- {
			closers[index]()
		}
	}
	idle := func() {}

	_, observer, failure := connector.base.provision(scope)
	if failure != nil {
		return integration.Installer{}, idle, failure
	}
	closers = append(closers, func() { session{observer: observer}.shutdown(scope) })

	recorder, store, failure := connector.base.recorder(scope, observer)
	if failure != nil {
		release()
		return integration.Installer{}, idle, failure
	}
	closers = append(closers, func() { _ = store.Close() })

	desk, failure := wiring.New(wiring.Options{Tracer: observer.Tracer(), Meter: observer.Meter()})
	if failure != nil {
		release()
		return integration.Installer{}, idle, failure
	}
	monitor, failure := telemetry.New(telemetry.Options{Tracer: observer.Tracer(), Meter: observer.Meter()})
	if failure != nil {
		release()
		return integration.Installer{}, idle, failure
	}
	installer, failure := integration.New(integration.Options{
		Resolver: binary.New(),
		Store:    desk,
		Recorder: recorder,
		Monitor:  monitor,
		Logger:   observer.Diagnostics(),
	})
	if failure != nil {
		release()
		return integration.Installer{}, idle, failure
	}
	return installer, release, nil
}
