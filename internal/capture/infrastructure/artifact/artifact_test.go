package artifact_test

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

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/artifact"
)

type fixture struct {
	test *testing.T
}

type setup struct {
	root    string
	ceiling int64
}

func (fixture fixture) make(setup setup) artifact.Store {
	fixture.test.Helper()
	store, failure := artifact.New(artifact.Options{
		Root:    setup.root,
		Ceiling: setup.ceiling,
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Tracer:  tracenoop.NewTracerProvider().Tracer("test"),
		Meter:   metricnoop.NewMeterProvider().Meter("test"),
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return store
}

func TestPutGet(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.make(setup{root: test.TempDir()})
	key, failure := store.Put(context.Background(), bytes.NewReader([]byte("hello")))
	if failure != nil {
		test.Fatal(failure)
	}
	payload, failure := store.Get(context.Background(), key)
	if failure != nil {
		test.Fatal(failure)
	}
	if string(payload) != "hello" {
		test.Fatalf("payload = %q", payload)
	}
}

func TestDedup(test *testing.T) {
	test.Parallel()

	root := test.TempDir()
	store := fixture{test: test}.make(setup{root: root})
	scope := context.Background()

	first, failure := store.Put(scope, bytes.NewReader([]byte("same")))
	if failure != nil {
		test.Fatal(failure)
	}
	second, failure := store.Put(scope, bytes.NewReader([]byte("same")))
	if failure != nil {
		test.Fatal(failure)
	}
	if first.String() != second.String() {
		test.Fatalf("identical content produced %q and %q", first.String(), second.String())
	}
	count := 0
	_ = filepath.WalkDir(filepath.Join(root, "objects"), func(_ string, entry os.DirEntry, _ error) error {
		if entry != nil && !entry.IsDir() {
			count++
		}
		return nil
	})
	if count != 1 {
		test.Fatalf("dedup left %d object files", count)
	}
}

func TestIntegrity(test *testing.T) {
	test.Parallel()

	root := test.TempDir()
	store := fixture{test: test}.make(setup{root: root})
	key, failure := store.Put(context.Background(), bytes.NewReader([]byte("data")))
	if failure != nil {
		test.Fatal(failure)
	}
	path := filepath.Join(root, "objects", key.String()[:2], key.String())
	if failure := os.WriteFile(path, []byte("tampered"), 0o600); failure != nil {
		test.Fatal(failure)
	}
	var integrity artifact.Integrity
	if _, failure := store.Get(context.Background(), key); !errors.As(failure, &integrity) {
		test.Fatalf("tampering was not detected: %v", failure)
	}
}

func TestOversize(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.make(setup{root: test.TempDir(), ceiling: 4})
	var oversize artifact.Oversize
	if _, failure := store.Put(context.Background(), bytes.NewReader([]byte("toolong"))); !errors.As(failure, &oversize) {
		test.Fatalf("oversize artifact was not rejected: %v", failure)
	}
}

func TestAbsent(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.make(setup{root: test.TempDir()})
	key, failure := digest.New(strings.Repeat("a", 64))
	if failure != nil {
		test.Fatal(failure)
	}
	var absent artifact.Absent
	if _, failure := store.Get(context.Background(), key); !errors.As(failure, &absent) {
		test.Fatalf("a missing artifact was not reported absent: %v", failure)
	}
}

func TestSweep(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.make(setup{root: test.TempDir()})
	orphan := filepath.Join(store.Temp(), "artifact-orphan")
	if failure := os.WriteFile(orphan, []byte("x"), 0o600); failure != nil {
		test.Fatal(failure)
	}
	if failure := store.Sweep(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := os.Stat(orphan); !errors.Is(failure, os.ErrNotExist) {
		test.Fatal("an orphaned temp was not swept")
	}
}
