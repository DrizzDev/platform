package connect_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/artifact"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/sqlite"
	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/telemetry"
)

// desk is a fake configuration store: it reports programmed presence and wired state, and counts the changes the
// installer asks it to make, so the flow can be exercised without touching a real file.
type desk struct {
	present     map[agent.Kind]bool
	wired       map[agent.Kind]bool
	captured    map[agent.Kind]bool
	fault       error
	connects    int
	disconnects int
	captures    int
	uncaptures  int
}

func (desk *desk) Detect(target agent.Agent) (bool, error) {
	return desk.present[target.Kind()], nil
}

func (desk *desk) Wired(target agent.Agent) (bool, error) {
	return desk.wired[target.Kind()], nil
}

func (desk *desk) Connect(context.Context, connect.Task) error {
	if desk.fault != nil {
		return desk.fault
	}
	desk.connects++
	return nil
}

func (desk *desk) Disconnect(context.Context, agent.Agent) error {
	if desk.fault != nil {
		return desk.fault
	}
	desk.disconnects++
	return nil
}

func (desk *desk) Captures(target agent.Agent) (bool, error) {
	return desk.captured[target.Kind()], nil
}

func (desk *desk) Capture(context.Context, connect.Task) error {
	if desk.fault != nil {
		return desk.fault
	}
	desk.captures++
	return nil
}

func (desk *desk) Uncapture(context.Context, agent.Agent) error {
	if desk.fault != nil {
		return desk.fault
	}
	desk.uncaptures++
	return nil
}

// pin is a fake resolver that yields a fixed executable path, or a failure.
type pin struct {
	path  string
	fault error
}

func (pin pin) Locate() (string, error) {
	return pin.path, pin.fault
}

type clock struct{}

func (clock) Now() time.Time {
	return time.Unix(1_000_000, 0)
}

type kit struct {
	test *testing.T
}

// setup is one installer under test: the fake store and resolver it is built over.
type setup struct {
	station  *desk
	resolver pin
}

// query asks for one agent's state in a report.
type query struct {
	report connect.Report
	kind   agent.Kind
}

func (kit kit) build(rig setup) (connect.Installer, sqlite.Store) {
	kit.test.Helper()
	dir := kit.test.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tracer := tracenoop.NewTracerProvider().Tracer("test")
	meter := metricnoop.NewMeterProvider().Meter("test")

	store, failure := sqlite.New(context.Background(), sqlite.Options{
		Path: filepath.Join(dir, "capture.db"), Logger: logger, Tracer: tracer, Meter: meter,
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	kit.test.Cleanup(func() { _ = store.Close() })
	vault, failure := artifact.New(artifact.Options{Root: dir, Logger: logger, Tracer: tracer, Meter: meter})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	made, failure := recording.New(recording.Options{
		Writer: store, Sink: vault, Keeper: store, Clock: clock{}, Logger: logger, Lease: time.Minute,
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	monitor, failure := telemetry.New(telemetry.Options{Tracer: tracer, Meter: meter})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	installer, failure := connect.New(connect.Options{
		Resolver: rig.resolver, Store: rig.station, Recorder: made, Monitor: monitor, Logger: logger,
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return installer, store
}

func (kit kit) state(ask query) connect.State {
	kit.test.Helper()
	for _, outcome := range ask.report.Outcomes() {
		if outcome.Kind() == ask.kind {
			return outcome.State()
		}
	}
	kit.test.Fatalf("no outcome for agent %q", ask.kind)
	return ""
}

func TestEnableConnectsDetectedAgents(test *testing.T) {
	test.Parallel()
	station := &desk{present: map[agent.Kind]bool{"claude-code": true, "codex": true}}
	kit := kit{test: test}
	installer, store := kit.build(setup{station: station, resolver: pin{path: "/opt/drizz"}})

	report, failure := installer.Enable(context.Background(), connect.Selection{All: true})
	if failure != nil {
		test.Fatal(failure)
	}
	if got := kit.state(query{report: report, kind: "claude-code"}); got != connect.Connected {
		test.Fatalf("claude-code = %s, want CONNECTED", got)
	}
	if got := kit.state(query{report: report, kind: "gemini"}); got != connect.Missing {
		test.Fatalf("gemini = %s, want MISSING", got)
	}
	if station.connects != 2 {
		test.Fatalf("connected %d agents, want 2", station.connects)
	}

	pending, failure := store.Pending(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if len(pending) != 2 {
		test.Fatalf("recorded %d installer actions, want 2", len(pending))
	}
}

func TestEnableUpdatesExistingEntry(test *testing.T) {
	test.Parallel()
	station := &desk{
		present: map[agent.Kind]bool{"claude-code": true},
		wired:   map[agent.Kind]bool{"claude-code": true},
	}
	kit := kit{test: test}
	installer, _ := kit.build(setup{station: station, resolver: pin{path: "/opt/drizz"}})

	report, failure := installer.Enable(context.Background(), connect.Selection{Kind: "claude-code"})
	if failure != nil {
		test.Fatal(failure)
	}
	if got := kit.state(query{report: report, kind: "claude-code"}); got != connect.Updated {
		test.Fatalf("claude-code = %s, want UPDATED", got)
	}
}

func TestEnableReportsFailurePerAgent(test *testing.T) {
	test.Parallel()
	station := &desk{present: map[agent.Kind]bool{"claude-code": true}, fault: connect.Malformed{}}
	kit := kit{test: test}
	installer, _ := kit.build(setup{station: station, resolver: pin{path: "/opt/drizz"}})

	report, failure := installer.Enable(context.Background(), connect.Selection{Kind: "claude-code"})
	if failure != nil {
		test.Fatal(failure)
	}
	if got := kit.state(query{report: report, kind: "claude-code"}); got != connect.Failed {
		test.Fatalf("claude-code = %s, want FAILED", got)
	}
}

func TestEnableAbortsWhenExecutableUnresolved(test *testing.T) {
	test.Parallel()
	station := &desk{present: map[agent.Kind]bool{"claude-code": true}}
	kit := kit{test: test}
	installer, _ := kit.build(setup{station: station, resolver: pin{fault: errors.New("no path")}})

	if _, failure := installer.Enable(context.Background(), connect.Selection{All: true}); failure == nil {
		test.Fatal("Enable must fail when the executable cannot be resolved")
	}
	if station.connects != 0 {
		test.Fatal("no agent may be touched when the executable cannot be resolved")
	}
}

func TestDisableRemovesEntry(test *testing.T) {
	test.Parallel()
	station := &desk{present: map[agent.Kind]bool{"codex": true}}
	kit := kit{test: test}
	installer, _ := kit.build(setup{station: station, resolver: pin{path: "/opt/drizz"}})

	report := installer.Disable(context.Background(), connect.Selection{Kind: "codex"})
	if got := kit.state(query{report: report, kind: "codex"}); got != connect.Removed {
		test.Fatalf("codex = %s, want REMOVED", got)
	}
	if station.disconnects != 1 {
		test.Fatalf("disconnected %d agents, want 1", station.disconnects)
	}
}

func TestSurveyReportsStatus(test *testing.T) {
	test.Parallel()
	station := &desk{
		present: map[agent.Kind]bool{"claude-code": true, "codex": true},
		wired:   map[agent.Kind]bool{"claude-code": true},
	}
	kit := kit{test: test}
	installer, _ := kit.build(setup{station: station, resolver: pin{path: "/opt/drizz"}})

	report := installer.Survey(context.Background())
	if got := kit.state(query{report: report, kind: "claude-code"}); got != connect.Connected {
		test.Fatalf("claude-code = %s, want CONNECTED", got)
	}
	if got := kit.state(query{report: report, kind: "codex"}); got != connect.Ready {
		test.Fatalf("codex = %s, want READY", got)
	}
	if got := kit.state(query{report: report, kind: "gemini"}); got != connect.Missing {
		test.Fatalf("gemini = %s, want MISSING", got)
	}
}
