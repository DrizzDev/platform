package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/retraction"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

func TestRetract(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))
	scope := context.Background()
	handle := fixture.handle("session_123")
	now := time.Unix(1000, 0)

	if _, failure := store.Publish(scope, fixture.notice(publication.Input{
		Session: handle, Expected: epoch.Epoch(0), Key: "key-1", Revision: 1, Moment: now,
	})); failure != nil {
		test.Fatal(failure)
	}

	request, failure := retraction.New(retraction.Input{Session: handle, Moment: now})
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := store.Retract(scope, request); failure != nil {
		test.Fatal(failure)
	}

	var absent sqlite.Absent
	if _, failure := store.Head(scope, handle); !errors.As(failure, &absent) {
		test.Fatalf("pointer was not cleared: %v", failure)
	}
	value, failure := store.Epoch(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if value != epoch.Epoch(2) {
		test.Fatalf("epoch = %d, want 2 after publish + retract", value)
	}
	backlog, failure := store.Backlog(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if backlog != 1 {
		test.Fatalf("active key was not queued for cleanup: %d", backlog)
	}
}

func TestRetractAbsent(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))
	request, failure := retraction.New(retraction.Input{Session: fixture.handle("session_404"), Moment: time.Unix(1000, 0)})
	if failure != nil {
		test.Fatal(failure)
	}

	var absent sqlite.Absent
	if failure := store.Retract(context.Background(), request); !errors.As(failure, &absent) {
		test.Fatalf("retract on an empty session should be absent: %v", failure)
	}
}
