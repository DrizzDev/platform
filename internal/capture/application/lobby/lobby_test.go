package lobby_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/application/lobby"
	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
)

var moment = time.Unix(10000, 0)

type register struct {
	held    []journal.Held
	evicted []int64
	next    int64
}

func (register *register) Admit(_ context.Context, observation journal.Observation) error {
	register.next++
	register.held = append(register.held, journal.Held{Reference: register.next, Observation: observation})
	return nil
}

func (register *register) Waiting(_ context.Context, window journal.Window) ([]journal.Held, error) {
	var live []journal.Held
	for _, item := range register.held {
		if !item.Observation.Moment.Before(window.Cutoff) {
			live = append(live, item)
		}
	}
	return live, nil
}

func (register *register) Evict(_ context.Context, references []int64) error {
	for _, reference := range references {
		register.evicted = append(register.evicted, reference)
		var kept []journal.Held
		for _, item := range register.held {
			if item.Reference != reference {
				kept = append(kept, item)
			}
		}
		register.held = kept
	}
	return nil
}

func (register *register) Expire(_ context.Context, cutoff time.Time) error {
	var kept []journal.Held
	for _, item := range register.held {
		if !item.Observation.Moment.Before(cutoff) {
			kept = append(kept, item)
		}
	}
	register.held = kept
	return nil
}

type clock struct{}

func (clock) Now() time.Time { return moment }

type spec struct {
	session     string
	fingerprint string
	moment      time.Time
	ordinal     int64
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) identity(value string) identifier.Identifier {
	fixture.test.Helper()
	made, failure := identifier.New(value)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) fingerprint(text string) digest.Digest {
	fixture.test.Helper()
	if text == "" {
		return digest.Digest{}
	}
	sum := sha256.Sum256([]byte(text))
	made, failure := digest.New(hex.EncodeToString(sum[:]))
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) crest(session string) bearings.Bearings {
	fixture.test.Helper()
	if session == "" {
		return bearings.New(bearings.Input{})
	}
	return bearings.New(bearings.Input{Session: fixture.identity(session)})
}

func (fixture fixture) observation(spec spec) journal.Observation {
	fixture.test.Helper()
	return journal.Observation{
		Bearings:    fixture.crest(spec.session),
		Fingerprint: fixture.fingerprint(spec.fingerprint),
		Moment:      spec.moment,
		Ordinal:     spec.ordinal,
	}
}

func (fixture fixture) call(spec spec) lobby.Call {
	fixture.test.Helper()
	return lobby.Call{
		Bearings:    fixture.crest(spec.session),
		Fingerprint: fixture.fingerprint(spec.fingerprint),
		Moment:      spec.moment,
		Ordinal:     spec.ordinal,
	}
}

func (fixture fixture) lobby(desk *register) lobby.Lobby {
	fixture.test.Helper()
	made, failure := lobby.New(lobby.Options{Register: desk, Clock: clock{}})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestExactClaim(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	desk := &register{}
	hall := kit.lobby(desk)
	scope := context.Background()

	if failure := hall.Observe(scope, kit.observation(spec{session: "s-1", moment: moment.Add(-time.Second), ordinal: 1})); failure != nil {
		test.Fatal(failure)
	}
	claimed, failure := hall.Activate(scope, kit.call(spec{session: "s-1", moment: moment, ordinal: 2}))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(claimed) != 1 || claimed[0].Mark != mark.Exact {
		test.Fatalf("claimed = %d, want one exact", len(claimed))
	}
}

func TestInferredClaim(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	desk := &register{}
	hall := kit.lobby(desk)
	scope := context.Background()

	if failure := hall.Observe(scope, kit.observation(spec{fingerprint: "input", moment: moment.Add(-time.Second), ordinal: 1})); failure != nil {
		test.Fatal(failure)
	}
	claimed, failure := hall.Activate(scope, kit.call(spec{fingerprint: "input", moment: moment, ordinal: 2}))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(claimed) != 1 || claimed[0].Mark != mark.Inferred {
		test.Fatalf("claimed = %d mark = %v, want one inferred", len(claimed), claimed)
	}
}

func TestNoClaim(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	desk := &register{}
	hall := kit.lobby(desk)
	scope := context.Background()

	if failure := hall.Observe(scope, kit.observation(spec{fingerprint: "one", moment: moment.Add(-time.Second), ordinal: 1})); failure != nil {
		test.Fatal(failure)
	}
	claimed, failure := hall.Activate(scope, kit.call(spec{fingerprint: "two", moment: moment, ordinal: 2}))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(claimed) != 0 {
		test.Fatalf("claimed = %d, want none", len(claimed))
	}
}

func TestClaimsEvery(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	desk := &register{}
	hall := kit.lobby(desk)
	scope := context.Background()

	for index := range 2 {
		observation := kit.observation(spec{session: "s-1", moment: moment.Add(-time.Second), ordinal: int64(index)})
		if failure := hall.Observe(scope, observation); failure != nil {
			test.Fatal(failure)
		}
	}
	claimed, failure := hall.Activate(scope, kit.call(spec{session: "s-1", moment: moment, ordinal: 9}))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(claimed) != 2 {
		test.Fatalf("claimed = %d, want both matching observations", len(claimed))
	}
}

func TestActivateEvicts(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	desk := &register{}
	hall := kit.lobby(desk)
	scope := context.Background()

	if failure := hall.Observe(scope, kit.observation(spec{session: "s-1", moment: moment.Add(-time.Second), ordinal: 1})); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := hall.Activate(scope, kit.call(spec{session: "s-1", moment: moment, ordinal: 2})); failure != nil {
		test.Fatal(failure)
	}
	again, failure := hall.Activate(scope, kit.call(spec{session: "s-1", moment: moment, ordinal: 3}))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(again) != 0 {
		test.Fatal("a claimed observation must not be activated twice")
	}
}

func TestExpiredNeverClaimed(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	desk := &register{}
	hall := kit.lobby(desk)
	scope := context.Background()

	if failure := hall.Observe(scope, kit.observation(spec{session: "s-1", moment: time.Unix(1, 0), ordinal: 1})); failure != nil {
		test.Fatal(failure)
	}
	claimed, failure := hall.Activate(scope, kit.call(spec{session: "s-1", moment: moment, ordinal: 2}))
	if failure != nil {
		test.Fatal(failure)
	}
	if len(claimed) != 0 || len(desk.held) != 0 {
		test.Fatal("an expired observation must be dropped, never claimed")
	}
}
