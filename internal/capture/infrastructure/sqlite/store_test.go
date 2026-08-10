package sqlite_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"sync"
	"testing"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/sqlite"
)

type fixture struct {
	test *testing.T
}

type seed struct {
	subject  trace.Trace
	hop      string
	sequence int64
}

func (fixture fixture) build(path string) (sqlite.Store, error) {
	return sqlite.New(context.Background(), sqlite.Options{
		Path:   path,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Meter:  metricnoop.NewMeterProvider().Meter("test"),
	})
}

func (fixture fixture) open(path string) sqlite.Store {
	fixture.test.Helper()
	store, failure := fixture.build(path)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	fixture.test.Cleanup(func() { _ = store.Close() })
	return store
}

func (fixture fixture) subject() trace.Trace {
	fixture.test.Helper()
	value, failure := trace.New("01HEXECUTION")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return value
}

func (fixture fixture) entry(seed seed) journal.Entry {
	fixture.test.Helper()
	here, failure := span.New(seed.hop)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	link, failure := correlation.New(correlation.Input{Trace: seed.subject, Span: here, Sequence: seed.sequence, Mark: mark.Exact})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	entry, failure := journal.New(journal.Input{
		Correlation: link,
		Origin:      origin.Capability,
		Fidelity:    fidelity.Exact,
		Category:    category.Tool,
		Payload:     []byte("payload"),
		State:       journal.Pending,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return entry
}

func TestAppendRead(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	subject := kit.subject()
	scope := context.Background()

	for index := range int64(3) {
		if failure := store.Append(scope, kit.entry(seed{subject: subject, hop: fmt.Sprintf("hop-%d", index), sequence: index})); failure != nil {
			test.Fatal(failure)
		}
	}
	entries, failure := store.Read(scope, subject)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(entries) != 3 {
		test.Fatalf("read %d entries", len(entries))
	}
	for index, entry := range entries {
		if entry.Correlation().Sequence() != int64(index) {
			test.Fatalf("entry %d out of order: sequence %d", index, entry.Correlation().Sequence())
		}
	}
}

func TestRollback(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	subject := kit.subject()
	scope, cancel := context.WithCancel(context.Background())
	cancel()

	if failure := store.Append(scope, kit.entry(seed{subject: subject, hop: "hop-0"})); !errors.Is(failure, context.Canceled) {
		test.Fatalf("cancellation did not propagate: %v", failure)
	}
	entries, failure := store.Read(context.Background(), subject)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(entries) != 0 {
		test.Fatalf("a cancelled append persisted %d entries", len(entries))
	}
}

func TestContention(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	path := filepath.Join(test.TempDir(), "capture.db")
	first := kit.open(path)
	second := kit.open(path)
	subject := kit.subject()

	var group sync.WaitGroup
	group.Add(2)
	go func() {
		defer group.Done()
		if failure := first.Append(context.Background(), kit.entry(seed{subject: subject, hop: "hop-0"})); failure != nil {
			test.Errorf("first append: %v", failure)
		}
	}()
	go func() {
		defer group.Done()
		if failure := second.Append(context.Background(), kit.entry(seed{subject: subject, hop: "hop-1", sequence: 1})); failure != nil {
			test.Errorf("second append: %v", failure)
		}
	}()
	group.Wait()

	entries, failure := first.Read(context.Background(), subject)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(entries) != 2 {
		test.Fatalf("contention lost a write: %d entries", len(entries))
	}
}

func TestCorruptRow(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	subject := kit.subject()

	insert := "INSERT INTO journal (trace, span, parent, sequence, mark, origin, fidelity, category, payload, state, stamped) " +
		"VALUES ('" + subject.String() + "', 'hop-0', '', 0, 'EXACT', 'GARBAGE', 'EXACT', 'TOOL', x'00', 'PENDING', 0)"
	if failure := store.Exec(insert); failure != nil {
		test.Fatal(failure)
	}
	var corrupt sqlite.Corrupt
	if _, failure := store.Read(context.Background(), subject); !errors.As(failure, &corrupt) {
		test.Fatalf("a corrupt row was not isolated: %v", failure)
	}
}
