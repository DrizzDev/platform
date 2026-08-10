package recording_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/artifact"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/sqlite"
)

type fixture struct {
	test *testing.T
}

func (fixture fixture) build() (recording.Recorder, sqlite.Store, artifact.Store) {
	fixture.test.Helper()
	dir := fixture.test.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tracer := tracenoop.NewTracerProvider().Tracer("test")
	meter := metricnoop.NewMeterProvider().Meter("test")

	store, failure := sqlite.New(context.Background(), sqlite.Options{
		Path: filepath.Join(dir, "capture.db"), Logger: logger, Tracer: tracer, Meter: meter,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	fixture.test.Cleanup(func() { _ = store.Close() })

	vault, failure := artifact.New(artifact.Options{Root: dir, Logger: logger, Tracer: tracer, Meter: meter})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	recorder, failure := recording.New(recording.Options{
		Writer: store, Sink: vault, Keeper: store, Clock: clock{}, Logger: logger, Lease: time.Minute,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return recorder, store, vault
}

type clock struct{}

func (clock) Now() time.Time { return time.Unix(1_000_000, 0) }

func TestRecord(test *testing.T) {
	test.Parallel()

	recorder, store, vault := fixture{test: test}.build()
	execution, failure := recorder.Begin()
	if failure != nil {
		test.Fatal(failure)
	}
	scope := context.Background()
	execution.Record(scope, recording.Note{
		Origin:   origin.Capability,
		Fidelity: fidelity.Exact,
		Category: category.Screen,
		Payload:  []byte("1080x2400"),
		Artifact: []byte("screenshot-bytes"),
	})

	entries, failure := store.Read(scope, execution.Trace())
	if failure != nil {
		test.Fatal(failure)
	}
	if len(entries) != 1 {
		test.Fatalf("recorded %d entries", len(entries))
	}
	reference := entries[0].Artifact()
	if reference.Empty() {
		test.Fatal("the entry carries no artifact reference")
	}
	payload, failure := vault.Get(scope, reference)
	if failure != nil {
		test.Fatal(failure)
	}
	if string(payload) != "screenshot-bytes" {
		test.Fatalf("artifact = %q", payload)
	}
}

func TestLeaseProtectsRun(test *testing.T) {
	test.Parallel()

	recorder, store, _ := fixture{test: test}.build()
	execution, failure := recorder.Begin()
	if failure != nil {
		test.Fatal(failure)
	}
	scope := context.Background()
	execution.Record(scope, recording.Note{Origin: origin.Capability, Fidelity: fidelity.Exact, Category: category.Log})

	claims, failure := store.Leases(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(claims) != 1 || claims[0].Trace.String() != execution.Trace().String() {
		test.Fatalf("recording did not lease the running trace, claims = %d", len(claims))
	}
}

type breaker struct{}

func (breaker) Append(context.Context, journal.Entry) error {
	return errors.New("store is gone")
}

func (breaker) Put(context.Context, io.Reader) (digest.Digest, error) {
	return digest.Digest{}, errors.New("store is gone")
}

func (breaker) Lease(context.Context, journal.Claim) error {
	return errors.New("store is gone")
}

func TestObservational(test *testing.T) {
	test.Parallel()

	var log bytes.Buffer
	recorder, failure := recording.New(recording.Options{
		Writer: breaker{},
		Sink:   breaker{},
		Keeper: breaker{},
		Clock:  clock{},
		Logger: slog.New(slog.NewJSONHandler(&log, nil)),
		Lease:  time.Minute,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	execution, failure := recorder.Begin()
	if failure != nil {
		test.Fatal(failure)
	}
	execution.Record(context.Background(), recording.Note{
		Origin:   origin.Capability,
		Fidelity: fidelity.Exact,
		Category: category.Screen,
		Artifact: []byte("shot"),
	})
	if !strings.Contains(log.String(), "capture.record.dropped") {
		test.Fatal("a failed write was not logged and dropped")
	}
}
