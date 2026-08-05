//go:build fencing

package identity_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/result"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

type contender struct {
	path string
	key  string
}

func (contender contender) open() (sqlite.Store, error) {
	return sqlite.New(context.Background(), sqlite.Options{
		Path:   contender.path,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Meter:  metricnoop.NewMeterProvider().Meter("test"),
	})
}

func (contender contender) run() int {
	store, failure := contender.open()
	if failure != nil {
		return 2
	}
	defer func() { _ = store.Close() }()

	handle, failure := session.New("session_123")
	if failure != nil {
		return 2
	}
	notice, failure := publication.New(publication.Input{
		Session:  handle,
		Expected: epoch.Epoch(0),
		Key:      contender.key,
		Revision: 1,
		Moment:   time.Unix(1000, 0),
	})
	if failure != nil {
		return 2
	}
	outcome, failure := store.Publish(context.Background(), notice)
	switch {
	case failure != nil:
		return 2
	case outcome == result.Published:
		return 0
	default:
		return 3
	}
}

type contest struct {
	test *testing.T
}

func (contest contest) spawn(path string) int {
	const workers = 5
	var group sync.WaitGroup
	codes := make(chan int, workers)
	for index := range workers {
		group.Add(1)
		go func(number int) {
			defer group.Done()
			command := exec.CommandContext(context.Background(), os.Args[0], "-test.run=^TestContention$")
			command.Env = append(os.Environ(), "FENCE_DATABASE="+path, "FENCE_KEY=key-"+strconv.Itoa(number))
			code := 0
			if failure := command.Run(); failure != nil {
				var exit *exec.ExitError
				if errors.As(failure, &exit) {
					code = exit.ExitCode()
				} else {
					code = -1
				}
			}
			codes <- code
		}(index)
	}
	group.Wait()
	close(codes)

	published := 0
	for code := range codes {
		switch code {
		case 0:
			published++
		case 3:
		default:
			contest.test.Fatalf("a worker failed with code %d", code)
		}
	}
	return published
}

func TestContention(test *testing.T) {
	if path := os.Getenv("FENCE_DATABASE"); path != "" {
		os.Exit(contender{path: path, key: os.Getenv("FENCE_KEY")}.run())
	}

	path := filepath.Join(test.TempDir(), "identity.db")
	warm, failure := (contender{path: path}).open()
	if failure != nil {
		test.Fatal(failure)
	}
	defer func() { _ = warm.Close() }()

	if published := (contest{test: test}).spawn(path); published != 1 {
		test.Fatalf("expected exactly one publish to win across processes, got %d", published)
	}
	value, failure := warm.Epoch(context.Background())
	if failure != nil {
		test.Fatal(failure)
	}
	if value != epoch.Epoch(1) {
		test.Fatalf("final epoch after the race = %d, want 1", value)
	}
}
