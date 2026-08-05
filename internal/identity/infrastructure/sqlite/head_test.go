package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

func TestHead(test *testing.T) {
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

	head, failure := store.Head(scope, handle)
	if failure != nil {
		test.Fatal(failure)
	}
	if head.Key() != "key-1" || head.Revision() != 1 || head.Epoch() != epoch.Epoch(1) {
		test.Fatalf("head = %+v", head)
	}
}

func TestHeadQuiet(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store, log := fixture.watch(filepath.Join(test.TempDir(), "identity.db"))

	var absent sqlite.Absent
	if _, failure := store.Head(context.Background(), fixture.handle("session_404")); !errors.As(failure, &absent) {
		test.Fatalf("expected absent: %v", failure)
	}
	if strings.Contains(log.String(), "identity.store.failed") {
		test.Fatalf("an absent head logged a fault:\n%s", log.String())
	}
}

func TestHeadAbsent(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))

	var absent sqlite.Absent
	if _, failure := store.Head(context.Background(), fixture.handle("session_404")); !errors.As(failure, &absent) {
		test.Fatalf("an absent head was not reported: %v", failure)
	}
}
