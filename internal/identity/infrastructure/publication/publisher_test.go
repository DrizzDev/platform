package publication_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/publication"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/vault"
)

type entropy struct{}

func (entropy) Bytes(size int) ([]byte, error) {
	return bytes.Repeat([]byte{9}, size), nil
}

type locker struct {
	entries map[string][]byte
	fail    bool
}

func (locker *locker) Read(_ context.Context, key string) ([]byte, error) {
	if locker.fail {
		return nil, errors.New("vault locked")
	}
	value, present := locker.entries[key]
	if !present {
		return nil, vault.Missing{}
	}
	return value, nil
}

func (locker *locker) Write(_ context.Context, entry vault.Entry) error {
	if locker.fail {
		return errors.New("vault locked")
	}
	locker.entries[entry.Key] = entry.Value
	return nil
}

func (locker *locker) Delete(_ context.Context, key string) error {
	delete(locker.entries, key)
	return nil
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) store() sqlite.Store {
	fixture.test.Helper()
	made, failure := sqlite.New(context.Background(), sqlite.Options{
		Path:   filepath.Join(fixture.test.TempDir(), "identity.db"),
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Meter:  metricnoop.NewMeterProvider().Meter("test"),
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	fixture.test.Cleanup(func() { _ = made.Close() })
	return made
}

func (fixture fixture) safe(inner vault.Store) vault.Vault {
	fixture.test.Helper()
	made, failure := vault.New(vault.Options{Store: inner})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) publisher(options publication.Options) publication.Publisher {
	fixture.test.Helper()
	made, failure := publication.New(options)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) handle() session.Session {
	fixture.test.Helper()
	made, failure := session.New("LOCAL")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) candidate(who string) login.Candidate {
	fixture.test.Helper()
	account, failure := subject.New(who)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return login.Candidate{
		Token: login.Token{
			Issuer:  "https://issuer.example/",
			Client:  "native",
			Subject: account,
			Method:  method.Browser,
			Refresh: []byte("refresh-token-bytes"),
			Issued:  time.Unix(1000, 0),
			Expiry:  time.Unix(2000, 0),
		},
		Moment: time.Unix(1500, 0),
	}
}

func TestPublish(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.store()
	safe := fixture.safe(&locker{entries: map[string][]byte{}})
	handle := fixture.handle()
	publisher := fixture.publisher(publication.Options{Vault: safe, Ledger: store, Random: entropy{}, Session: handle})
	scope := context.Background()

	receipt, failure := publisher.Publish(scope, fixture.candidate("google-oauth2|first"))
	if failure != nil {
		test.Fatal(failure)
	}
	if receipt.Subject.String() != "google-oauth2|first" || receipt.Session.String() != "LOCAL" {
		test.Fatalf("receipt = %+v", receipt)
	}
	head, failure := store.Head(scope, handle)
	if failure != nil {
		test.Fatal(failure)
	}
	if head.Revision() != 1 || !strings.HasPrefix(head.Key(), "LOCAL#1#") {
		test.Fatalf("head = %+v", head)
	}
}

func TestSupersede(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.store()
	safe := fixture.safe(&locker{entries: map[string][]byte{}})
	handle := fixture.handle()
	publisher := fixture.publisher(publication.Options{Vault: safe, Ledger: store, Random: entropy{}, Session: handle})
	scope := context.Background()

	if _, failure := publisher.Publish(scope, fixture.candidate("google-oauth2|first")); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := publisher.Publish(scope, fixture.candidate("google-oauth2|first")); failure != nil {
		test.Fatal(failure)
	}
	head, failure := store.Head(scope, handle)
	if failure != nil {
		test.Fatal(failure)
	}
	if head.Revision() != 2 || !strings.HasPrefix(head.Key(), "LOCAL#2#") {
		test.Fatalf("head = %+v", head)
	}
	backlog, failure := store.Backlog(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if backlog != 1 {
		test.Fatalf("superseded credential not queued for cleanup: %d", backlog)
	}
}

func TestConflict(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.store()
	safe := fixture.safe(&locker{entries: map[string][]byte{}})
	handle := fixture.handle()
	publisher := fixture.publisher(publication.Options{Vault: safe, Ledger: store, Random: entropy{}, Session: handle})
	scope := context.Background()

	if _, failure := publisher.Publish(scope, fixture.candidate("google-oauth2|first")); failure != nil {
		test.Fatal(failure)
	}
	_, failure := publisher.Publish(scope, fixture.candidate("google-oauth2|second"))
	var conflict login.Conflict
	if !errors.As(failure, &conflict) {
		test.Fatalf("a different account was allowed to switch: %v", failure)
	}
}

func (fixture fixture) saturate(store sqlite.Store) {
	fixture.test.Helper()
	for index := range 16 {
		record, failure := cleanup.New(cleanup.Input{
			Key: "orphan#" + strconv.Itoa(index), Reason: cleanup.Superseded, State: cleanup.Pending,
			Next: time.Unix(1, 0), Deadline: time.Unix(9000, 0), Created: time.Unix(1, 0),
		})
		if failure != nil {
			fixture.test.Fatal(failure)
		}
		if failure := store.Enqueue(context.Background(), record); failure != nil {
			fixture.test.Fatal(failure)
		}
	}
}

func TestSaturated(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.store()
	safe := fixture.safe(&locker{entries: map[string][]byte{}})
	handle := fixture.handle()
	publisher := fixture.publisher(publication.Options{Vault: safe, Ledger: store, Random: entropy{}, Session: handle})
	fixture.saturate(store)

	_, failure := publisher.Publish(context.Background(), fixture.candidate("google-oauth2|first"))
	var storage login.Storage
	if !errors.As(failure, &storage) {
		test.Fatalf("a saturated backlog did not block the new credential: %v", failure)
	}
}

func TestStorage(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.store()
	safe := fixture.safe(&locker{entries: map[string][]byte{}, fail: true})
	handle := fixture.handle()
	publisher := fixture.publisher(publication.Options{Vault: safe, Ledger: store, Random: entropy{}, Session: handle})

	_, failure := publisher.Publish(context.Background(), fixture.candidate("google-oauth2|first"))
	var storage login.Storage
	if !errors.As(failure, &storage) {
		test.Fatalf("a vault failure was not reported as storage: %v", failure)
	}
}
