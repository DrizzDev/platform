package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/result"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

type query struct {
	test  *testing.T
	store sqlite.Store
}

func (query query) value(statement string) string {
	query.test.Helper()
	value, failure := query.store.Query(statement)
	if failure != nil {
		query.test.Fatal(failure)
	}
	return value
}

func (query query) expect(checks map[string]string) {
	query.test.Helper()
	for statement, want := range checks {
		if value := query.value(statement); value != want {
			query.test.Fatalf("%s = %q, want %q", statement, value, want)
		}
	}
}

func TestPublish(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))
	scope := context.Background()
	handle := fixture.handle("session_123")
	now := time.Unix(1000, 0)
	query := query{test: test, store: store}

	first, failure := store.Publish(scope, fixture.notice(publication.Input{Session: handle, Expected: epoch.Epoch(0), Key: "key-1", Revision: 1, Moment: now}))
	if failure != nil {
		test.Fatal(failure)
	}
	refresh, failure := store.Publish(scope, fixture.notice(publication.Input{Session: handle, Expected: epoch.Epoch(1), Key: "key-2", Revision: 2, Moment: now}))
	if failure != nil {
		test.Fatal(failure)
	}
	if first != result.Published || refresh != result.Published {
		test.Fatalf("publish results = %q %q", first, refresh)
	}
	query.expect(map[string]string{
		"SELECT value FROM epoch WHERE id = 1":                                                              "2",
		"SELECT name FROM pointer WHERE session = 'session_123'":                                            "key-2",
		"SELECT count(*) FROM cleanup WHERE name = 'key-1' AND reason = 'SUPERSEDED' AND state = 'PENDING'": "1",
	})
}

func TestReject(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.open(filepath.Join(test.TempDir(), "identity.db"))
	scope := context.Background()
	handle := fixture.handle("session_123")
	now := time.Unix(1000, 0)
	query := query{test: test, store: store}

	if _, failure := store.Publish(scope, fixture.notice(publication.Input{Session: handle, Expected: epoch.Epoch(0), Key: "key-1", Revision: 1, Moment: now})); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := store.Publish(scope, fixture.notice(publication.Input{Session: handle, Expected: epoch.Epoch(1), Key: "key-2", Revision: 2, Moment: now})); failure != nil {
		test.Fatal(failure)
	}

	stale, failure := store.Publish(scope, fixture.notice(publication.Input{Session: handle, Expected: epoch.Epoch(1), Key: "key-3", Revision: 3, Moment: now}))
	if failure != nil {
		test.Fatal(failure)
	}
	if stale != result.Rejected {
		test.Fatalf("stale publish = %q", stale)
	}
	query.expect(map[string]string{
		"SELECT value FROM epoch WHERE id = 1":                                                            "2",
		"SELECT name FROM pointer WHERE session = 'session_123'":                                          "key-2",
		"SELECT count(*) FROM cleanup WHERE name = 'key-3' AND reason = 'REJECTED' AND state = 'PENDING'": "1",
	})
}
