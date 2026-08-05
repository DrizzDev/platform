package vault_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/vault"
)

const window = 2560

type store struct {
	entries map[string][]byte
}

func (store *store) Read(_ context.Context, key string) ([]byte, error) {
	value, found := store.entries[key]
	if !found {
		return nil, vault.Missing{}
	}
	return value, nil
}

func (store *store) Write(_ context.Context, entry vault.Entry) error {
	store.entries[entry.Key] = entry.Value
	return nil
}

func (store *store) Delete(_ context.Context, key string) error {
	delete(store.entries, key)
	return nil
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) record(refresh []byte) credential.Record {
	fixture.test.Helper()
	account, failure := subject.New("google-oauth2|110000000000000000000")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	handle, failure := session.New("123e4567-e89b-12d3-a456-426614174000")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	record, failure := credential.New(credential.Input{
		Issuer:   "https://issuer.example/",
		Client:   "aBcDeFgHiJkLmNoPqRsTuVwXyZ012345",
		Handle:   "handle_1234567890",
		Subject:  account,
		Session:  handle,
		Method:   method.Browser,
		Refresh:  refresh,
		Issued:   time.Unix(1000, 0),
		Expiry:   time.Unix(2000, 0),
		Revision: 1,
		Schema:   1,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return record
}

func (fixture fixture) vault(store *store) vault.Vault {
	fixture.test.Helper()
	made, failure := vault.New(vault.Options{Store: store})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestRoundTrip(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := &store{entries: map[string][]byte{}}
	safe := fixture.vault(store)
	record := fixture.record(bytes.Repeat([]byte("r"), 256))

	if failure := safe.Write(context.Background(), record); failure != nil {
		test.Fatal(failure)
	}
	restored, failure := safe.Read(context.Background(), record.Key())
	if failure != nil {
		test.Fatal(failure)
	}
	if restored.Session().String() != record.Session().String() ||
		!bytes.Equal(restored.Refresh(), record.Refresh()) ||
		restored.Method() != method.Browser {
		test.Fatalf("restored = %+v", restored)
	}

	stored := store.entries[string(record.Key())]
	if len(stored) >= window {
		test.Fatalf("stored record %d bytes exceeds the Windows blob limit %d", len(stored), window)
	}
	test.Logf("realistic record encodes to %d bytes (Windows blob limit %d)", len(stored), window)
}

func TestBudget(test *testing.T) {
	test.Parallel()

	account, failure := subject.New(strings.Repeat("s", 256))
	if failure != nil {
		test.Fatal(failure)
	}
	handle, failure := session.New(strings.Repeat("h", 256))
	if failure != nil {
		test.Fatal(failure)
	}
	record, failure := credential.New(credential.Input{
		Issuer:   strings.Repeat("i", 512),
		Client:   strings.Repeat("c", 512),
		Handle:   strings.Repeat("h", 43),
		Subject:  account,
		Session:  handle,
		Method:   method.Browser,
		Refresh:  bytes.Repeat([]byte("r"), 512),
		Issued:   time.Unix(1000, 0),
		Expiry:   time.Unix(2000, 0),
		Revision: 1,
		Schema:   1,
	})
	if failure != nil {
		test.Fatal(failure)
	}

	store := &store{entries: map[string][]byte{}}
	safe := fixture{test: test}.vault(store)
	if failure := safe.Write(context.Background(), record); failure != nil {
		test.Fatalf("a maximal valid record was not storable: %v", failure)
	}
	test.Logf("maximal record encodes to %d bytes (limit %d)", len(store.entries[string(record.Key())]), window)
}

func TestOverflow(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := &store{entries: map[string][]byte{}}
	safe := fixture.vault(store)
	key := fixture.record(bytes.Repeat([]byte("r"), 8)).Key()
	store.entries[string(key)] = make([]byte, window+1)

	if _, failure := safe.Read(context.Background(), key); failure == nil {
		test.Fatal("an oversized stored blob was accepted")
	}
}

func TestMissing(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	safe := fixture.vault(&store{entries: map[string][]byte{}})
	key := fixture.record(bytes.Repeat([]byte("r"), 8)).Key()

	_, absent := safe.Read(context.Background(), key)
	var missing vault.Missing
	if !errors.As(absent, &missing) {
		test.Fatalf("missing credential = %v", absent)
	}
}

func TestDependencies(test *testing.T) {
	test.Parallel()

	if _, failure := vault.New(vault.Options{}); failure == nil {
		test.Fatal("a vault without a store was accepted")
	}
}
