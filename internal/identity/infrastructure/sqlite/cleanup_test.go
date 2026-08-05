package sqlite_test

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/deferral"
)

type maker struct {
	test *testing.T
}

func (maker maker) record(key string) cleanup.Record {
	maker.test.Helper()
	record, failure := cleanup.New(cleanup.Input{
		Key:      key,
		Reason:   cleanup.Superseded,
		State:    cleanup.Pending,
		Attempts: 0,
		Next:     time.Unix(1000, 0),
		Deadline: time.Unix(9000, 0),
		Created:  time.Unix(500, 0),
	})
	if failure != nil {
		maker.test.Fatal(failure)
	}
	return record
}

func (maker maker) postpone(key string) deferral.Deferral {
	maker.test.Helper()
	postponed, failure := deferral.New(deferral.Input{Key: key, Next: time.Unix(1000, 0)})
	if failure != nil {
		maker.test.Fatal(failure)
	}
	return postponed
}

func TestCleanup(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.open(filepath.Join(test.TempDir(), "identity.db"))
	scope := context.Background()
	maker := maker{test: test}

	for index := range 5 {
		if failure := store.Enqueue(scope, maker.record("key-"+strconv.Itoa(index))); failure != nil {
			test.Fatal(failure)
		}
	}

	pending, failure := store.Pending(scope, time.Unix(1500, 0))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(pending) != 4 {
		test.Fatalf("pending = %d, want the reconcile limit of 4", len(pending))
	}

	early, failure := store.Pending(scope, time.Unix(500, 0))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(early) != 0 {
		test.Fatalf("records due later were returned: %d", len(early))
	}

	if failure := store.Acknowledge(scope, "key-0"); failure != nil {
		test.Fatal(failure)
	}
	count, failure := store.Backlog(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if count != 4 {
		test.Fatalf("backlog = %d after acknowledge", count)
	}
}

func TestPlan(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.open(filepath.Join(test.TempDir(), "identity.db"))
	plan, failure := store.Plan("SELECT name, reason, state, attempts, next, deadline, created FROM cleanup " +
		"WHERE state = 'PENDING' AND next <= 1500 ORDER BY created LIMIT 4")
	if failure != nil {
		test.Fatal(failure)
	}
	if !strings.Contains(plan, "USING INDEX due") {
		test.Fatalf("the reconcile query does not use the index:\n%s", plan)
	}
}

func TestBacklog(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.open(filepath.Join(test.TempDir(), "identity.db"))
	scope := context.Background()
	maker := maker{test: test}

	for index := range 1000 {
		if failure := store.Enqueue(scope, maker.record("key-"+strconv.Itoa(index))); failure != nil {
			test.Fatal(failure)
		}
	}
	pending, failure := store.Pending(scope, time.Unix(1500, 0))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(pending) != 4 {
		test.Fatalf("a large backlog broke the reconcile batch bound: %d", len(pending))
	}
	count, failure := store.Backlog(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if count != 1000 {
		test.Fatalf("backlog = %d, want 1000", count)
	}
}

func TestDefer(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.open(filepath.Join(test.TempDir(), "identity.db"))
	scope := context.Background()
	maker := maker{test: test}
	query := query{test: test, store: store}
	postponed := maker.postpone("key-0")

	if failure := store.Enqueue(scope, maker.record("key-0")); failure != nil {
		test.Fatal(failure)
	}
	for range 4 {
		if failure := store.Defer(scope, postponed); failure != nil {
			test.Fatal(failure)
		}
	}
	if query.value("SELECT state FROM cleanup WHERE name = 'key-0'") != "PENDING" {
		test.Fatal("a record blocked before exhausting its retries")
	}

	if failure := store.Defer(scope, postponed); failure != nil {
		test.Fatal(failure)
	}
	if query.value("SELECT state FROM cleanup WHERE name = 'key-0'") != "BLOCKED" {
		test.Fatal("an exhausted record was not blocked")
	}
	pending, failure := store.Pending(scope, time.Unix(2000, 0))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(pending) != 0 {
		test.Fatalf("a blocked record was returned as pending: %d", len(pending))
	}
}
