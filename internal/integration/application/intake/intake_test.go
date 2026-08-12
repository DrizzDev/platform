package intake_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/artifact"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/sqlite"
	"github.com/DrizzDev/platform/internal/integration/application/intake"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/telemetry"
)

type clock struct{}

func (clock) Now() time.Time {
	return time.Unix(1_000_000, 0)
}

type kit struct {
	test *testing.T
}

func (kit kit) build() (intake.Intake, sqlite.Store) {
	test := kit.test
	test.Helper()
	dir := test.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tracer := tracenoop.NewTracerProvider().Tracer("test")
	meter := metricnoop.NewMeterProvider().Meter("test")

	store, failure := sqlite.New(context.Background(), sqlite.Options{
		Path: filepath.Join(dir, "capture.db"), Logger: logger, Tracer: tracer, Meter: meter,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	test.Cleanup(func() { _ = store.Close() })
	vault, failure := artifact.New(artifact.Options{Root: dir, Logger: logger, Tracer: tracer, Meter: meter})
	if failure != nil {
		test.Fatal(failure)
	}
	made, failure := recording.New(recording.Options{
		Writer: store, Sink: vault, Keeper: store, Clock: clock{}, Logger: logger, Lease: time.Minute,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	monitor, failure := telemetry.New(telemetry.Options{Tracer: tracer, Meter: meter})
	if failure != nil {
		test.Fatal(failure)
	}
	receiver, failure := intake.New(intake.Options{Recorder: made, Monitor: monitor, Logger: logger})
	if failure != nil {
		test.Fatal(failure)
	}
	return receiver, store
}

func TestRecordPromptIsHostPrompt(test *testing.T) {
	test.Parallel()
	receiver, store := kit{test: test}.build()
	scope := context.Background()

	receiver.Record(scope, intake.Event{Agent: "claude-code", Slot: agent.Prompt, Text: "open settings"})

	pending, failure := store.Pending(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(pending) != 1 {
		test.Fatalf("recorded %d events, want 1", len(pending))
	}
	if pending[0].Origin() != origin.Host {
		test.Fatalf("origin = %s, want HOST", pending[0].Origin())
	}
	if pending[0].Category() != category.Prompt {
		test.Fatalf("category = %s, want PROMPT", pending[0].Category())
	}
}

func TestRecordTurnIsResponse(test *testing.T) {
	test.Parallel()
	receiver, store := kit{test: test}.build()
	scope := context.Background()

	receiver.Record(scope, intake.Event{Agent: "codex", Slot: agent.Turn, Text: "done"})

	pending, failure := store.Pending(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(pending) != 1 || pending[0].Category() != category.Response {
		test.Fatalf("want one RESPONSE record, got %d", len(pending))
	}
}

// probe is a fake monitor that captures every string the intake hands to telemetry, so a canary can prove no captured
// content is among them.
type probe struct {
	seen []string
}

func (probe *probe) Watch(scope context.Context, operation string) (context.Context, func(string)) {
	probe.seen = append(probe.seen, operation)
	return scope, func(outcome string) { probe.seen = append(probe.seen, outcome) }
}

// witnessed is one instrumented intake under observation: the receiver, the artifact store the text should reach, the
// probe monitor and the log buffer the text must never reach.
type witnessed struct {
	receiver intake.Intake
	store    sqlite.Store
	vault    artifact.Store
	watcher  *probe
	log      *bytes.Buffer
}

func (kit kit) witness() witnessed {
	kit.test.Helper()
	dir := kit.test.TempDir()
	log := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(log, nil))
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
	watcher := &probe{}
	receiver, failure := intake.New(intake.Options{Recorder: made, Monitor: watcher, Logger: logger})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return witnessed{receiver: receiver, store: store, vault: vault, watcher: watcher, log: log}
}

// TestTelemetryOmitsPromptText is the privacy canary: a captured prompt's text must reach the durable capture store
// but never any telemetry sink. It records a prompt with a unique secret, then proves the secret is in the stored
// artifact yet absent from every log line and every string the monitor received.
func TestTelemetryOmitsPromptText(test *testing.T) {
	test.Parallel()
	const secret = "canary-prompt-7f3a9c-do-not-leak"

	rig := kit{test: test}.witness()
	scope := context.Background()
	rig.receiver.Record(scope, intake.Event{Agent: "claude-code", Slot: agent.Prompt, Text: secret})

	pending, failure := rig.store.Pending(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(pending) != 1 {
		test.Fatalf("recorded %d events, want 1", len(pending))
	}
	blob, failure := rig.vault.Get(scope, pending[0].Artifact())
	if failure != nil {
		test.Fatal(failure)
	}
	if string(blob) != secret {
		test.Fatalf("the prompt text must be preserved in the capture artifact, got %q", blob)
	}
	if strings.Contains(rig.log.String(), secret) {
		test.Fatal("prompt text leaked into the telemetry log")
	}
	for _, value := range rig.watcher.seen {
		if strings.Contains(value, secret) {
			test.Fatalf("prompt text leaked into a telemetry label: %q", value)
		}
	}
}
