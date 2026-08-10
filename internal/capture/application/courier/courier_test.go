package courier_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/capture/application/courier"
	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

// book is a fake journal: it yields pending entries and, unless frozen (to model a lost ack), removes an entry when it
// is advanced.
type book struct {
	pending  []journal.Entry
	advanced []string
	frozen   bool
}

func (book *book) Pending(context.Context) ([]journal.Entry, error) {
	return book.pending, nil
}

func (book *book) Advance(_ context.Context, transition journal.Transition) error {
	span := transition.Entry.Correlation().Span().String()
	book.advanced = append(book.advanced, span)
	if book.frozen {
		return nil
	}
	kept := make([]journal.Entry, 0, len(book.pending))
	for _, entry := range book.pending {
		if entry.Correlation().Span().String() != span {
			kept = append(kept, entry)
		}
	}
	book.pending = kept
	return nil
}

type store struct {
	blobs map[string][]byte
}

func (store store) Get(_ context.Context, key digest.Digest) ([]byte, error) {
	blob, found := store.blobs[key.String()]
	if !found {
		return nil, errors.New("artifact missing")
	}
	return blob, nil
}

// cloud is a fake uploader: it dedups by identity, counts every call, and can be told to fail transiently or reject.
type cloud struct {
	seen   map[string]bool
	order  []string
	faults int
	calls  int
	reject bool
}

func (cloud *cloud) Blob(_ context.Context, cargo courier.Cargo) error {
	return cloud.receive("artifact:" + cargo.Digest.String())
}

func (cloud *cloud) Record(_ context.Context, entry journal.Entry) error {
	return cloud.receive("entry:" + entry.Correlation().Span().String())
}

func (cloud *cloud) receive(identity string) error {
	cloud.calls++
	if cloud.reject {
		return courier.Rejected{}
	}
	if cloud.faults > 0 {
		cloud.faults--
		return errors.New("temporary cloud outage")
	}
	if cloud.seen == nil {
		cloud.seen = map[string]bool{}
	}
	if cloud.seen[identity] {
		return nil
	}
	cloud.seen[identity] = true
	cloud.order = append(cloud.order, identity)
	return nil
}

type fixture struct {
	test  *testing.T
	vault store
}

type spec struct {
	hop  string
	blob []byte
}

func (fixture fixture) entry(spec spec) journal.Entry {
	fixture.test.Helper()
	thread, failure := trace.New("01HEXECUTION")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	point, failure := span.New(spec.hop)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	link, failure := correlation.New(correlation.Input{Trace: thread, Span: point, Mark: mark.Exact})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	reference := digest.Digest{}
	if len(spec.blob) > 0 {
		sum := sha256.Sum256(spec.blob)
		reference, failure = digest.New(hex.EncodeToString(sum[:]))
		if failure != nil {
			fixture.test.Fatal(failure)
		}
		fixture.vault.blobs[reference.String()] = spec.blob
	}
	entry, failure := journal.New(journal.Input{
		Correlation: link,
		Origin:      origin.Capability,
		Fidelity:    fidelity.Exact,
		Category:    category.Tool,
		Artifact:    reference,
		State:       journal.Pending,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return entry
}

type rig struct {
	ledger *book
	sky    *cloud
}

func (fixture fixture) build(rig rig) courier.Courier {
	fixture.test.Helper()
	carrier, failure := courier.New(courier.Options{Ledger: rig.ledger, Vault: fixture.vault, Uploader: rig.sky})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return carrier
}

func TestDrain(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test, vault: store{blobs: map[string][]byte{}}}
	ledger := &book{pending: []journal.Entry{
		kit.entry(spec{hop: "a"}),
		kit.entry(spec{hop: "b", blob: []byte("shot")}),
	}}
	sky := &cloud{}
	if failure := kit.build(rig{ledger: ledger, sky: sky}).Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if len(ledger.pending) != 0 {
		test.Fatalf("%d entries left pending", len(ledger.pending))
	}
	if len(ledger.advanced) != 2 {
		test.Fatalf("advanced %d entries", len(ledger.advanced))
	}
	if len(sky.seen) != 3 || !sky.seen["entry:a"] || !sky.seen["entry:b"] {
		test.Fatalf("uploaded set = %v", sky.seen)
	}
}

func TestOrder(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test, vault: store{blobs: map[string][]byte{}}}
	ledger := &book{pending: []journal.Entry{kit.entry(spec{hop: "a", blob: []byte("shot")})}}
	sky := &cloud{}
	if failure := kit.build(rig{ledger: ledger, sky: sky}).Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if len(sky.order) != 2 || !strings.HasPrefix(sky.order[0], "artifact:") || !strings.HasPrefix(sky.order[1], "entry:") {
		test.Fatalf("upload order = %v", sky.order)
	}
}

func TestIdempotent(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test, vault: store{blobs: map[string][]byte{}}}
	ledger := &book{frozen: true, pending: []journal.Entry{kit.entry(spec{hop: "a", blob: []byte("shot")})}}
	sky := &cloud{}
	carrier := kit.build(rig{ledger: ledger, sky: sky})
	if failure := carrier.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if failure := carrier.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if len(sky.seen) != 2 {
		test.Fatalf("re-send was not deduped: %d distinct effects", len(sky.seen))
	}
	if sky.calls <= 2 {
		test.Fatalf("the second run did not re-send: %d calls", sky.calls)
	}
}

func TestRetry(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test, vault: store{blobs: map[string][]byte{}}}
	ledger := &book{pending: []journal.Entry{kit.entry(spec{hop: "a"})}}
	sky := &cloud{faults: 2}
	if failure := kit.build(rig{ledger: ledger, sky: sky}).Run(context.Background()); failure != nil {
		test.Fatalf("a transient failure was not retried: %v", failure)
	}
	if !sky.seen["entry:a"] {
		test.Fatal("the entry was not delivered after retry")
	}
}

func TestReject(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test, vault: store{blobs: map[string][]byte{}}}
	ledger := &book{pending: []journal.Entry{kit.entry(spec{hop: "a"})}}
	sky := &cloud{reject: true}
	failure := kit.build(rig{ledger: ledger, sky: sky}).Run(context.Background())
	var rejected courier.Rejected
	if !errors.As(failure, &rejected) {
		test.Fatalf("failure = %v", failure)
	}
	if sky.calls != 1 {
		test.Fatalf("a rejected upload was retried: %d calls", sky.calls)
	}
	if len(ledger.advanced) != 0 {
		test.Fatal("a rejected entry was marked synced")
	}
}
