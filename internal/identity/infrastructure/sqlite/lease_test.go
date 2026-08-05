package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/hold"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

func TestLease(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))
	scope := context.Background()
	handle := fixture.handle("session_123")
	now := time.Unix(1000, 0)
	span := 30 * time.Second

	held, failure := store.Acquire(scope, fixture.grip(hold.Input{Session: handle, Owner: "one", Moment: now, Window: span}))
	if failure != nil {
		test.Fatal(failure)
	}
	if held.Owner() != "one" {
		test.Fatalf("owner = %q", held.Owner())
	}

	var contention sqlite.Contention
	_, contended := store.Acquire(scope, fixture.grip(hold.Input{Session: handle, Owner: "two", Moment: now.Add(time.Second), Window: span}))
	if !errors.As(contended, &contention) {
		test.Fatalf("a held lease was taken: %v", contended)
	}

	if _, failure := store.Acquire(scope, fixture.grip(hold.Input{Session: handle, Owner: "two", Moment: now.Add(time.Minute), Window: span})); failure != nil {
		test.Fatalf("an expired lease was not re-acquirable: %v", failure)
	}

	if failure := store.Renew(scope, fixture.grip(hold.Input{Session: handle, Owner: "one", Moment: now.Add(time.Minute), Window: span})); !errors.As(failure, &contention) {
		test.Fatalf("a non-owner renewed the lease: %v", failure)
	}
	if failure := store.Renew(scope, fixture.grip(hold.Input{Session: handle, Owner: "two", Moment: now.Add(90 * time.Second), Window: span})); failure != nil {
		test.Fatalf("the owner could not renew: %v", failure)
	}

	if failure := store.Release(scope, fixture.grip(hold.Input{Session: handle, Owner: "two", Moment: now, Window: span})); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := store.Acquire(scope, fixture.grip(hold.Input{Session: handle, Owner: "one", Moment: now.Add(2 * time.Minute), Window: span})); failure != nil {
		test.Fatalf("the lease was not free after release: %v", failure)
	}
}
