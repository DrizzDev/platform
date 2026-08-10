package sqlite_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
)

func (fixture fixture) reference(text string) digest.Digest {
	fixture.test.Helper()
	sum := sha256.Sum256([]byte(text))
	made, failure := digest.New(hex.EncodeToString(sum[:]))
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

type bound struct {
	seed      seed
	reference digest.Digest
}

func (fixture fixture) linked(bound bound) journal.Entry {
	fixture.test.Helper()
	here, failure := span.New(bound.seed.hop)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	link, failure := correlation.New(correlation.Input{
		Trace: bound.seed.subject, Span: here, Sequence: bound.seed.sequence, Mark: mark.Exact,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	entry, failure := journal.New(journal.Input{
		Correlation: link,
		Origin:      origin.Capability,
		Fidelity:    fidelity.Exact,
		Category:    category.Screen,
		Payload:     []byte("payload"),
		Artifact:    bound.reference,
		State:       journal.Pending,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return entry
}

func TestSettled(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	subject := kit.subject()
	scope := context.Background()

	first := kit.entry(seed{subject: subject, hop: "hop-0", sequence: 0})
	second := kit.entry(seed{subject: subject, hop: "hop-1", sequence: 1})
	pending := kit.entry(seed{subject: subject, hop: "hop-2", sequence: 2})
	for _, entry := range []journal.Entry{first, second, pending} {
		if failure := store.Append(scope, entry); failure != nil {
			test.Fatal(failure)
		}
	}
	for _, entry := range []journal.Entry{first, second} {
		if failure := store.Advance(scope, journal.Transition{Entry: entry, State: journal.Synced}); failure != nil {
			test.Fatal(failure)
		}
	}

	retained, failure := store.Settled(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(retained) != 2 {
		test.Fatalf("settled = %d, want 2 synced (un-synced excluded)", len(retained))
	}
	if retained[0].Span.String() != "hop-0" || retained[1].Span.String() != "hop-1" {
		test.Fatalf("settled order = %q, %q", retained[0].Span.String(), retained[1].Span.String())
	}
	if retained[0].Category != category.Tool {
		test.Fatalf("category = %q", retained[0].Category)
	}
}

func TestDiscard(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	subject := kit.subject()
	scope := context.Background()

	keep := kit.entry(seed{subject: subject, hop: "hop-keep", sequence: 0})
	drop := kit.entry(seed{subject: subject, hop: "hop-drop", sequence: 1})
	for _, entry := range []journal.Entry{keep, drop} {
		if failure := store.Append(scope, entry); failure != nil {
			test.Fatal(failure)
		}
	}
	if failure := store.Discard(scope, []span.Span{drop.Correlation().Span()}); failure != nil {
		test.Fatal(failure)
	}
	remaining, failure := store.Read(scope, subject)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(remaining) != 1 || remaining[0].Correlation().Span().String() != "hop-keep" {
		test.Fatalf("after discard = %d entries", len(remaining))
	}
}

func TestDigests(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	subject := kit.subject()
	scope := context.Background()

	shared := kit.reference("shared")
	other := kit.reference("other")
	entries := []journal.Entry{
		kit.linked(bound{seed: seed{subject: subject, hop: "hop-0", sequence: 0}, reference: shared}),
		kit.linked(bound{seed: seed{subject: subject, hop: "hop-1", sequence: 1}, reference: shared}),
		kit.linked(bound{seed: seed{subject: subject, hop: "hop-2", sequence: 2}, reference: other}),
		kit.entry(seed{subject: subject, hop: "hop-3", sequence: 3}),
	}
	for _, entry := range entries {
		if failure := store.Append(scope, entry); failure != nil {
			test.Fatal(failure)
		}
	}
	references, failure := store.Digests(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(references) != 2 {
		test.Fatalf("digests = %d, want 2 distinct (dedup shared, ignore empty)", len(references))
	}
}

func TestLease(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	subject := kit.subject()
	scope := context.Background()

	if failure := store.Lease(scope, journal.Claim{Trace: subject, Until: time.Unix(1000, 0)}); failure != nil {
		test.Fatal(failure)
	}
	if failure := store.Lease(scope, journal.Claim{Trace: subject, Until: time.Unix(5000, 0)}); failure != nil {
		test.Fatal(failure)
	}
	claims, failure := store.Leases(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if len(claims) != 1 {
		test.Fatalf("leases = %d, want 1 (a repeat lease extends, not duplicates)", len(claims))
	}
	if !claims[0].Until.Equal(time.Unix(5000, 0)) || claims[0].Trace.String() != subject.String() {
		test.Fatalf("claim = %s until %v", claims[0].Trace.String(), claims[0].Until)
	}
}
