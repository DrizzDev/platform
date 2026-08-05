package sqlite_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/attempt"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/hold"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

type fixture struct {
	test *testing.T
}

func (fixture fixture) build(path string) (sqlite.Store, error) {
	return sqlite.New(context.Background(), sqlite.Options{
		Path:   path,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Meter:  metricnoop.NewMeterProvider().Meter("test"),
	})
}

func (fixture fixture) watch(path string) (sqlite.Store, *bytes.Buffer) {
	fixture.test.Helper()
	var log bytes.Buffer
	store, failure := sqlite.New(context.Background(), sqlite.Options{
		Path:   path,
		Logger: slog.New(slog.NewJSONHandler(&log, nil)),
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Meter:  metricnoop.NewMeterProvider().Meter("test"),
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	fixture.test.Cleanup(func() { _ = store.Close() })
	return store, &log
}

func (fixture fixture) open(path string) sqlite.Store {
	fixture.test.Helper()
	store, failure := fixture.build(path)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	fixture.test.Cleanup(func() {
		_ = store.Close()
	})
	return store
}

func (fixture fixture) handle(value string) session.Session {
	fixture.test.Helper()
	handle, failure := session.New(value)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return handle
}

func (fixture fixture) notice(input publication.Input) publication.Publication {
	fixture.test.Helper()
	notice, failure := publication.New(input)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return notice
}

func (fixture fixture) trial(input attempt.Input) attempt.Attempt {
	fixture.test.Helper()
	trial, failure := attempt.New(input)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return trial
}

func (fixture fixture) mark(input marking.Input) marking.Marking {
	fixture.test.Helper()
	mark, failure := marking.New(input)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return mark
}

func (fixture fixture) grip(input hold.Input) hold.Hold {
	fixture.test.Helper()
	grip, failure := hold.New(input)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return grip
}

func TestMigrate(test *testing.T) {
	test.Parallel()

	path := filepath.Join(test.TempDir(), "identity.db")
	store := fixture{test: test}.open(path)

	checks := map[string]string{
		"PRAGMA user_version": "1",
		"PRAGMA journal_mode": "wal",
		"PRAGMA foreign_keys": "1",
		"SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name IN ('epoch', 'pointer', 'attempt', 'lease', 'cleanup')": "5",
		"SELECT value FROM epoch WHERE id = 1": "0",
	}
	for statement, expected := range checks {
		value, failure := store.Query(statement)
		if failure != nil {
			test.Fatalf("%s: %v", statement, failure)
		}
		if value != expected {
			test.Fatalf("%s = %q, want %q", statement, value, expected)
		}
	}
}

func TestIdempotent(test *testing.T) {
	test.Parallel()

	path := filepath.Join(test.TempDir(), "identity.db")
	fixture := fixture{test: test}
	fixture.open(path)

	reopened := fixture.open(path)
	value, failure := reopened.Query("PRAGMA user_version")
	if failure != nil {
		test.Fatal(failure)
	}
	if value != "1" {
		test.Fatalf("reopened version = %q", value)
	}
}

func TestCorrupt(test *testing.T) {
	test.Parallel()

	path := filepath.Join(test.TempDir(), "identity.db")
	if failure := os.WriteFile(path, []byte("this is not a database"), 0o600); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := (fixture{test: test}).build(path); failure == nil {
		test.Fatal("a corrupt database file was opened")
	}
}

func TestDependencies(test *testing.T) {
	test.Parallel()

	if _, failure := sqlite.New(context.Background(), sqlite.Options{}); failure == nil {
		test.Fatal("a database without options was accepted")
	}
}

func TestStorage(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.open(filepath.Join(test.TempDir(), "identity.db"))
	pages, failure := store.Query("PRAGMA page_count")
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := store.Exec("PRAGMA max_page_count = " + pages); failure != nil {
		test.Fatal(failure)
	}
	bulk := strings.Repeat("x", 300000)
	insert := "INSERT INTO cleanup (name, reason, state, attempts, next, deadline, created) VALUES ('" +
		bulk + "', 'LOGOUT', 'PENDING', 0, 0, 0, 0)"
	if failure := store.Exec(insert); failure == nil {
		test.Fatal("a write past the page ceiling was accepted")
	}
}

func TestRollback(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))
	scope, cancel := context.WithCancel(context.Background())
	cancel()

	notice := fixture.notice(publication.Input{
		Session: fixture.handle("session_123"), Expected: epoch.Epoch(0), Key: "key-1", Revision: 1, Moment: time.Unix(1000, 0),
	})
	if _, failure := store.Publish(scope, notice); !errors.Is(failure, context.Canceled) {
		test.Fatalf("cancellation did not propagate: %v", failure)
	}
	if value, _ := store.Query("SELECT value FROM epoch WHERE id = 1"); value != "0" {
		test.Fatalf("a cancelled publish advanced the epoch to %q", value)
	}
}
