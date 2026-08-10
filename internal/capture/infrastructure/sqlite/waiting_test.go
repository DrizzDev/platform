package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
)

type watch struct {
	moment   time.Time
	session  string
	payload  string
	category category.Category
	ordinal  int64
}

func (fixture fixture) identity(value string) identifier.Identifier {
	fixture.test.Helper()
	made, failure := identifier.New(value)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) observation(watch watch) journal.Observation {
	fixture.test.Helper()
	return journal.Observation{
		Bearings: bearings.New(bearings.Input{Session: fixture.identity(watch.session)}),
		Moment:   watch.moment,
		Ordinal:  watch.ordinal,
		Origin:   origin.Host,
		Fidelity: fidelity.Exact,
		Category: watch.category,
		Payload:  []byte(watch.payload),
	}
}

func TestAdmitWaiting(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	scope := context.Background()

	for _, spec := range []watch{
		{session: "s-1", moment: time.Unix(1000, 0), ordinal: 1, category: category.Prompt, payload: "one"},
		{session: "s-2", moment: time.Unix(1001, 0), ordinal: 2, category: category.Response, payload: "two"},
	} {
		if failure := store.Admit(scope, kit.observation(spec)); failure != nil {
			test.Fatal(failure)
		}
	}
	held, failure := store.Waiting(scope, journal.Window{Cutoff: time.Unix(0, 0), Capacity: 10})
	if failure != nil {
		test.Fatal(failure)
	}
	if len(held) != 2 || held[0].Observation.Ordinal != 2 {
		test.Fatalf("waiting = %d, want 2 newest-first", len(held))
	}
	if held[0].Observation.Bearings.Session().String() != "s-2" || string(held[0].Observation.Payload) != "two" {
		test.Fatal("observation did not round-trip through the store")
	}
}

func TestWaitingBounded(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	scope := context.Background()

	for index := range 5 {
		spec := watch{session: "s-1", moment: time.Unix(int64(1000+index), 0), ordinal: int64(index), category: category.Tool}
		if failure := store.Admit(scope, kit.observation(spec)); failure != nil {
			test.Fatal(failure)
		}
	}
	held, failure := store.Waiting(scope, journal.Window{Cutoff: time.Unix(0, 0), Capacity: 2})
	if failure != nil {
		test.Fatal(failure)
	}
	if len(held) != 2 {
		test.Fatalf("waiting = %d, want capacity-bounded to 2", len(held))
	}
}

func TestEvict(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	scope := context.Background()

	if failure := store.Admit(scope, kit.observation(watch{session: "s-1", moment: time.Unix(1000, 0), category: category.Tool})); failure != nil {
		test.Fatal(failure)
	}
	held, failure := store.Waiting(scope, journal.Window{Cutoff: time.Unix(0, 0), Capacity: 10})
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := store.Evict(scope, []int64{held[0].Reference}); failure != nil {
		test.Fatal(failure)
	}
	remaining, failure := store.Waiting(scope, journal.Window{Cutoff: time.Unix(0, 0), Capacity: 10})
	if failure != nil {
		test.Fatal(failure)
	}
	if len(remaining) != 0 {
		test.Fatalf("evicted observation still present, waiting = %d", len(remaining))
	}
}

func TestExpire(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	store := kit.open(filepath.Join(test.TempDir(), "capture.db"))
	scope := context.Background()

	stale := kit.observation(watch{session: "s-1", moment: time.Unix(100, 0), category: category.Tool})
	fresh := kit.observation(watch{session: "s-2", moment: time.Unix(1000, 0), category: category.Tool})
	for _, observation := range []journal.Observation{stale, fresh} {
		if failure := store.Admit(scope, observation); failure != nil {
			test.Fatal(failure)
		}
	}
	if failure := store.Expire(scope, time.Unix(500, 0)); failure != nil {
		test.Fatal(failure)
	}
	held, failure := store.Waiting(scope, journal.Window{Cutoff: time.Unix(0, 0), Capacity: 10})
	if failure != nil {
		test.Fatal(failure)
	}
	if len(held) != 1 || held[0].Observation.Bearings.Session().String() != "s-2" {
		test.Fatalf("expire kept the wrong rows, waiting = %d", len(held))
	}
}
