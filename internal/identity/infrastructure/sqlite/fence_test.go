package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/attempt"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

func TestFence(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))
	scope := context.Background()
	handle := fixture.handle("session_123")

	if _, failure := store.Publish(scope, fixture.notice(publication.Input{
		Session: handle, Expected: epoch.Epoch(0), Key: "key-1", Revision: 1, Moment: time.Unix(1000, 0),
	})); failure != nil {
		test.Fatal(failure)
	}

	active := fixture.mark(marking.Input{Session: handle, Attempt: fixture.trial(attempt.Input{Revision: 1, Epoch: epoch.Epoch(1)})})
	if failure := store.Fence(scope, active); failure != nil {
		test.Fatalf("the active revision could not be fenced: %v", failure)
	}

	var fenced sqlite.Fenced
	if failure := store.Fence(scope, active); !errors.As(failure, &fenced) {
		test.Fatalf("a revision was attempted twice: %v", failure)
	}

	stale := fixture.mark(marking.Input{Session: handle, Attempt: fixture.trial(attempt.Input{Revision: 2, Epoch: epoch.Epoch(1)})})
	var contention sqlite.Contention
	if failure := store.Fence(scope, stale); !errors.As(failure, &contention) {
		test.Fatalf("a stale revision was fenced: %v", failure)
	}
}

func TestFenceAbsent(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))
	mark := fixture.mark(marking.Input{
		Session: fixture.handle("session_123"),
		Attempt: fixture.trial(attempt.Input{Revision: 1, Epoch: epoch.Epoch(1)}),
	})

	var contention sqlite.Contention
	if failure := store.Fence(context.Background(), mark); !errors.As(failure, &contention) {
		test.Fatalf("fencing without an active credential was accepted: %v", failure)
	}
}
