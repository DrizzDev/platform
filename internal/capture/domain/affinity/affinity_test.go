package affinity_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/affinity"
	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
)

const window = 10 * time.Second

var moment = time.Unix(1000, 0)

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

func (fixture fixture) digest(text string) digest.Digest {
	fixture.test.Helper()
	sum := sha256.Sum256([]byte(text))
	made, failure := digest.New(hex.EncodeToString(sum[:]))
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) session(value string) bearings.Bearings {
	fixture.test.Helper()
	return bearings.New(bearings.Input{Session: fixture.identity(value)})
}

func TestExact(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	call := affinity.New(affinity.Input{Window: window, Bearings: kit.session("s-1"), Moment: moment, Ordinal: 2})
	candidate := affinity.New(affinity.Input{Bearings: kit.session("s-1"), Moment: moment, Ordinal: 1})

	link, matched := call.Match(candidate)
	if !matched || link != mark.Exact {
		test.Fatalf("a shared identifier must be exact, got %q matched=%v", link, matched)
	}
}

func TestInferredByInput(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	shared := kit.digest("normalized-input")
	call := affinity.New(affinity.Input{Window: window, Fingerprint: shared, Moment: moment, Ordinal: 2})
	candidate := affinity.New(affinity.Input{Fingerprint: shared, Moment: moment.Add(-time.Second), Ordinal: 1})

	link, matched := call.Match(candidate)
	if !matched || link != mark.Inferred {
		test.Fatalf("matching input within the window must be inferred, got %q matched=%v", link, matched)
	}
}

func TestInferredByProcess(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	process := kit.identity("pid-9")
	call := affinity.New(affinity.Input{Window: window, Process: process, Moment: moment, Ordinal: 2})
	candidate := affinity.New(affinity.Input{Process: process, Moment: moment.Add(-time.Second), Ordinal: 1})

	link, matched := call.Match(candidate)
	if !matched || link != mark.Inferred {
		test.Fatalf("same process within the window must be inferred, got %q matched=%v", link, matched)
	}
}

func TestBeyondWindow(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	shared := kit.digest("normalized-input")
	call := affinity.New(affinity.Input{Window: window, Fingerprint: shared, Moment: moment, Ordinal: 2})
	candidate := affinity.New(affinity.Input{Fingerprint: shared, Moment: moment.Add(-time.Hour), Ordinal: 1})

	if _, matched := call.Match(candidate); matched {
		test.Fatal("an observation outside the window must not match")
	}
}

func TestOutOfOrder(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	shared := kit.digest("normalized-input")
	call := affinity.New(affinity.Input{Window: window, Fingerprint: shared, Moment: moment, Ordinal: 1})
	candidate := affinity.New(affinity.Input{Fingerprint: shared, Moment: moment, Ordinal: 5})

	if _, matched := call.Match(candidate); matched {
		test.Fatal("an observation seen after the call must not infer a match")
	}
}

func TestMissing(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	call := affinity.New(affinity.Input{Window: window, Fingerprint: kit.digest("one"), Moment: moment, Ordinal: 2})
	candidate := affinity.New(affinity.Input{Fingerprint: kit.digest("two"), Moment: moment, Ordinal: 1})

	if _, matched := call.Match(candidate); matched {
		test.Fatal("no shared identifier, input, or process must be no match")
	}
}
