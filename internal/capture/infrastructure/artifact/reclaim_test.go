package artifact_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/infrastructure/artifact"
)

func TestPrune(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.make(setup{root: test.TempDir()})
	scope := context.Background()

	keep, failure := store.Put(scope, bytes.NewReader([]byte("keep")))
	if failure != nil {
		test.Fatal(failure)
	}
	drop, failure := store.Put(scope, bytes.NewReader([]byte("drop")))
	if failure != nil {
		test.Fatal(failure)
	}
	referenced := func(reference digest.Digest, _ time.Time) bool { return reference.String() == keep.String() }
	if failure := store.Prune(scope, referenced); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := store.Get(scope, keep); failure != nil {
		test.Fatalf("kept object gone: %v", failure)
	}
	if _, failure := store.Get(scope, drop); !errors.As(failure, &artifact.Absent{}) {
		test.Fatalf("unreferenced object survived prune: %v", failure)
	}
	if failure := store.Prune(scope, referenced); failure != nil {
		test.Fatalf("prune must be idempotent, got %v", failure)
	}
}

func TestPruneSparesYoung(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.make(setup{root: test.TempDir()})
	scope := context.Background()

	young, failure := store.Put(scope, bytes.NewReader([]byte("young")))
	if failure != nil {
		test.Fatal(failure)
	}
	spare := func(_ digest.Digest, modified time.Time) bool { return time.Since(modified) < time.Hour }
	if failure := store.Prune(scope, spare); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := store.Get(scope, young); failure != nil {
		test.Fatalf("young object was pruned despite the grace predicate: %v", failure)
	}
}

func TestFootprintArtifact(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.make(setup{root: test.TempDir()})
	scope := context.Background()

	if _, failure := store.Put(scope, bytes.NewReader([]byte("hello"))); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := store.Put(scope, bytes.NewReader([]byte("world!"))); failure != nil {
		test.Fatal(failure)
	}
	total, failure := store.Footprint(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	if total != 11 {
		test.Fatalf("footprint = %d, want 11 (5 + 6)", total)
	}
}

func TestGate(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.make(setup{root: test.TempDir()})
	scope := context.Background()

	if failure := store.Restrict(scope); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := store.Put(scope, bytes.NewReader([]byte("blocked"))); !errors.As(failure, &artifact.Saturated{}) {
		test.Fatalf("restricted Put = %v, want Saturated", failure)
	}
	if failure := store.Admit(scope); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := store.Put(scope, bytes.NewReader([]byte("allowed"))); failure != nil {
		test.Fatalf("after admit Put = %v, want success", failure)
	}
}
