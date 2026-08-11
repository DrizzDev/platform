package host

import (
	"context"
	"sync"
	"time"

	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/infrastructure/bridge"
	"github.com/DrizzDev/platform/internal/device/infrastructure/sidecar"
	"github.com/DrizzDev/platform/internal/transport/cli"
	"github.com/DrizzDev/platform/internal/transport/mcp"
)

// dispatch bounds one device request to the helper; a slow single call fails on its own without stalling the rest.
const dispatch = 30 * time.Second

var (
	_ cli.Perform = pilot{}
	_ mcp.Perform = pilot{}
)

// pilot is the host's device runtime, shared by both surfaces. It holds a station that assembles the compiled helper,
// the device connection, and the recorder once on first use and reuses them, so a long-lived agent connection keeps
// one helper running rather than restarting it per call. Until the helper is configured every capability reports the
// same typed, non-retryable message, and the same runtime works end to end once the helper is in place. Cloud
// capabilities will follow this identical shape.
type pilot struct {
	station *station
}

func (pilot pilot) Screenshot(scope context.Context, target operator.Target) (operator.Shot, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Shot{}, failure
	}
	return engine.Screenshot(scope, target)
}

func (pilot pilot) Devices(scope context.Context) (operator.Roster, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Roster{}, failure
	}
	return engine.Devices(scope)
}

// station assembles the device runtime once and holds it for reuse. Assembly is single-flight; its result — the
// operator, or the reason it is unavailable — is kept for the life of the connection, and close tears the runtime down.
type station struct {
	base   foundation
	once   sync.Once
	engine operator.Operator
	stop   func()
	fault  error
}

// operate returns the shared operator, assembling it once. The assembly is detached from the triggering request so the
// helper it starts is not torn down when that request completes.
func (station *station) operate(scope context.Context) (operator.Operator, error) {
	station.once.Do(func() {
		station.engine, station.stop, station.fault = station.base.assemble(context.WithoutCancel(scope))
	})
	return station.engine, station.fault
}

func (station *station) close() {
	if station.stop != nil {
		station.stop()
	}
}

// assemble builds the device operator — the compiled helper, the device connection, and the local recorder — and
// returns a release that shuts them down in reverse. Until the helper is configured the capability is unprepared,
// which the surfaces render as one clear message; installation resolves the helper here without touching the rest.
func (foundation foundation) assemble(scope context.Context) (operator.Operator, func(), error) {
	var closers []func()
	release := func() {
		for index := len(closers) - 1; index >= 0; index-- {
			closers[index]()
		}
	}
	idle := func() {}

	handle, ready := sidecar.New(foundation.environment).Locate()
	if !ready {
		return operator.Operator{}, idle, operator.Refusal{Code: outcome.Unprepared}
	}

	_, observer, failure := foundation.provision(scope)
	if failure != nil {
		return operator.Operator{}, idle, failure
	}
	closers = append(closers, func() { session{observer: observer}.shutdown(scope) })

	recorder, store, failure := foundation.recorder(scope, observer)
	if failure != nil {
		release()
		return operator.Operator{}, idle, failure
	}
	closers = append(closers, func() { _ = store.Close() })

	driver, failure := bridge.New(bridge.Options{Location: handle.Location, Digest: handle.Digest, Timeout: dispatch})
	if failure != nil {
		release()
		return operator.Operator{}, idle, failure
	}
	closers = append(closers, func() { _ = driver.Close() })

	flow, failure := control.New(control.Options{Bridge: driver})
	if failure != nil {
		release()
		return operator.Operator{}, idle, failure
	}

	engine, failure := operator.New(operator.Options{Flow: flow, Recorder: recorder, Logger: observer.Diagnostics()})
	if failure != nil {
		release()
		return operator.Operator{}, idle, failure
	}
	return engine, release, nil
}
