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

func (pilot pilot) Snapshot(scope context.Context, target operator.Target) (operator.Snapshot, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Snapshot{}, failure
	}
	return engine.Snapshot(scope, target)
}

func (pilot pilot) Hierarchy(scope context.Context, target operator.Target) (operator.Tree, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Tree{}, failure
	}
	return engine.Hierarchy(scope, target)
}

func (pilot pilot) Dimensions(scope context.Context, target operator.Target) (operator.Extent, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Extent{}, failure
	}
	return engine.Dimensions(scope, target)
}

func (pilot pilot) Tap(scope context.Context, contact operator.Contact) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Tap(scope, contact)
}

func (pilot pilot) Install(scope context.Context, target operator.Package) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Install(scope, target)
}

func (pilot pilot) Launch(scope context.Context, target operator.Application) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Launch(scope, target)
}

func (pilot pilot) Terminate(scope context.Context, target operator.Application) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Terminate(scope, target)
}

func (pilot pilot) Wipe(scope context.Context, target operator.Application) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Wipe(scope, target)
}

func (pilot pilot) Installed(scope context.Context, target operator.Target) (operator.Listing, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Listing{}, failure
	}
	return engine.Installed(scope, target)
}

func (pilot pilot) Running(scope context.Context, target operator.Target) (operator.Listing, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Listing{}, failure
	}
	return engine.Running(scope, target)
}

func (pilot pilot) Foreground(scope context.Context, target operator.Target) (operator.Report, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Report{}, failure
	}
	return engine.Foreground(scope, target)
}

func (pilot pilot) Url(scope context.Context, target operator.Target) (operator.Report, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Report{}, failure
	}
	return engine.Url(scope, target)
}

func (pilot pilot) Disk(scope context.Context, target operator.Target) (operator.Measure, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Measure{}, failure
	}
	return engine.Disk(scope, target)
}

func (pilot pilot) Images(scope context.Context) (operator.Images, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Images{}, failure
	}
	return engine.Images(scope)
}

func (pilot pilot) Boot(scope context.Context, target operator.Image) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Boot(scope, target)
}

func (pilot pilot) Pause(scope context.Context, target operator.Target) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Pause(scope, target)
}

func (pilot pilot) Resume(scope context.Context, target operator.Target) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Resume(scope, target)
}

func (pilot pilot) Swipe(scope context.Context, drag operator.Drag) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Swipe(scope, drag)
}

func (pilot pilot) Pinch(scope context.Context, squeeze operator.Squeeze) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Pinch(scope, squeeze)
}

func (pilot pilot) Press(scope context.Context, key operator.Key) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Press(scope, key)
}

func (pilot pilot) Type(scope context.Context, input operator.Entry) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Type(scope, input)
}

func (pilot pilot) Clear(scope context.Context, target operator.Target) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Clear(scope, target)
}

func (pilot pilot) Back(scope context.Context, target operator.Target) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Back(scope, target)
}

func (pilot pilot) Home(scope context.Context, target operator.Target) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Home(scope, target)
}

func (pilot pilot) Locate(scope context.Context, fix operator.Fix) (operator.Ack, error) {
	engine, failure := pilot.station.operate(scope)
	if failure != nil {
		return operator.Ack{}, failure
	}
	return engine.Locate(scope, fix)
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
