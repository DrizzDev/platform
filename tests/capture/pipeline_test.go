package capture_test

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/capture/application/courier"
	"github.com/DrizzDev/platform/internal/capture/application/janitor"
	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/application/lobby"
	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/catalogue"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/domain/policy"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/artifact"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/sqlite"
)

// cloud is a fake acknowledging every upload, so the courier can drive an end-to-end synchronization.
type cloud struct {
	records map[string]bool
}

func (cloud *cloud) Blob(context.Context, courier.Cargo) error { return nil }

func (cloud *cloud) Record(_ context.Context, entry journal.Entry) error {
	cloud.records[entry.Correlation().Span().String()] = true
	return nil
}

type bell struct{}

func (bell) Alert(context.Context, janitor.Pressure) error { return nil }

type clock struct{}

func (clock) Now() time.Time { return time.Now() }

// distant is far past any real record time, so retention has elapsed and the recorder's short-lived lease has lapsed
// when the janitor runs — making reclaim deterministic without sleeping.
type distant struct{}

func (distant) Now() time.Time { return time.Unix(1<<40, 0) }

type kit struct {
	test  *testing.T
	store sqlite.Store
	vault artifact.Store
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) assemble() kit {
	fixture.test.Helper()
	dir := fixture.test.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	meter := metricnoop.NewMeterProvider().Meter("test")
	tracer := tracenoop.NewTracerProvider().Tracer("test")

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
	return kit{test: fixture.test, store: store, vault: vault}
}

func (kit kit) begin() *recording.Execution {
	kit.test.Helper()
	recorder, failure := recording.New(recording.Options{
		Writer: kit.store, Sink: kit.vault, Keeper: kit.store, Clock: clock{},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), Lease: time.Nanosecond,
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	execution, failure := recorder.Begin()
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return execution
}

func (kit kit) recorded(execution *recording.Execution) []journal.Entry {
	kit.test.Helper()
	entries, failure := kit.store.Read(context.Background(), execution.Trace())
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return entries
}

func (kit kit) hall() lobby.Lobby {
	kit.test.Helper()
	made, failure := lobby.New(lobby.Options{Register: kit.store, Clock: clock{}})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return made
}

func (kit kit) sync(sky *cloud) {
	kit.test.Helper()
	post, failure := courier.New(courier.Options{Ledger: kit.store, Vault: kit.vault, Uploader: sky})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	if failure := post.Run(context.Background()); failure != nil {
		kit.test.Fatal(failure)
	}
}

func (kit kit) reclaim() janitor.Report {
	kit.test.Helper()
	rule, failure := policy.New(policy.Input{Limit: 1024, Retention: time.Millisecond, Upload: true})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	shelf, failure := catalogue.New(catalogue.Input{Policies: map[category.Category]policy.Policy{
		category.Prompt: rule,
		category.Tool:   rule,
	}})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	keeper, failure := janitor.New(janitor.Options{
		Archive: kit.store, Vault: kit.vault, Notifier: bell{}, Clock: distant{}, Catalogue: shelf,
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	report, failure := keeper.Run(context.Background())
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return report
}

func (kit kit) call(session string) lobby.Call {
	kit.test.Helper()
	return lobby.Call{Bearings: kit.crest(session), Moment: time.Now(), Ordinal: 9}
}

func (kit kit) observation(session string) journal.Observation {
	kit.test.Helper()
	return journal.Observation{
		Bearings: kit.crest(session),
		Moment:   time.Now().Add(-time.Second),
		Ordinal:  1,
		Origin:   origin.Host,
		Fidelity: fidelity.Exact,
		Category: category.Prompt,
		Payload:  []byte("prompt"),
	}
}

func (kit kit) crest(session string) bearings.Bearings {
	kit.test.Helper()
	made, failure := identifier.New(session)
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return bearings.New(bearings.Input{Session: made})
}

// TestPipeline drives one host observation through the whole capture stack: it is held in the pending window, claimed
// by a capability call, recorded into the execution, synchronized to the cloud, then reclaimed once acknowledged.
func TestPipeline(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}.assemble()
	scope := context.Background()

	hall := kit.hall()
	if failure := hall.Observe(scope, kit.observation("s-1")); failure != nil {
		test.Fatal(failure)
	}

	execution := kit.begin()
	claimed, failure := hall.Activate(scope, kit.call("s-1"))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(claimed) != 1 || claimed[0].Mark != mark.Exact {
		test.Fatalf("claimed = %d, want one exact host observation", len(claimed))
	}
	for _, taken := range claimed {
		execution.Record(scope, recording.Note{
			Origin: taken.Observation.Origin, Fidelity: taken.Observation.Fidelity,
			Category: taken.Observation.Category, Payload: taken.Observation.Payload, Mark: taken.Mark,
		})
	}
	execution.Record(scope, recording.Note{
		Origin: origin.Capability, Fidelity: fidelity.Exact, Category: category.Tool, Payload: []byte("result"),
	})

	recorded := kit.recorded(execution)
	if len(recorded) != 2 || recorded[0].Correlation().Mark() != mark.Exact {
		test.Fatalf("recorded %d entries, want the claimed observation and the capability step", len(recorded))
	}

	sky := &cloud{records: map[string]bool{}}
	kit.sync(sky)
	if len(sky.records) != 2 {
		test.Fatalf("cloud received %d records, want 2", len(sky.records))
	}

	if report := kit.reclaim(); report.Reclaimed != 2 {
		test.Fatalf("reclaimed %d, want the two acknowledged records", report.Reclaimed)
	}
	if remaining := kit.recorded(execution); len(remaining) != 0 {
		test.Fatalf("acknowledged records were not reclaimed, %d remain", len(remaining))
	}
}
