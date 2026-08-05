package marking_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/attempt"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

type fixture struct {
	test *testing.T
}

func (fixture fixture) handle(value string) session.Session {
	fixture.test.Helper()
	handle, failure := session.New(value)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return handle
}

func (fixture fixture) trial(input attempt.Input) attempt.Attempt {
	fixture.test.Helper()
	trial, failure := attempt.New(input)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return trial
}

func TestMarking(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	mark, failure := marking.New(marking.Input{
		Session: fixture.handle("session_123"),
		Attempt: fixture.trial(attempt.Input{Revision: 1, Epoch: epoch.Epoch(1)}),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if mark.Attempt().Revision() != 1 || mark.Session().String() != "session_123" {
		test.Fatalf("marking = %+v", mark)
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	if _, failure := marking.New(marking.Input{Attempt: fixture.trial(attempt.Input{Revision: 1, Epoch: epoch.Epoch(1)})}); failure == nil {
		test.Fatal("a marking without a session was accepted")
	}
}
