package hold_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/hold"
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

func TestHold(test *testing.T) {
	test.Parallel()

	moment := time.Unix(1000, 0)
	held, failure := hold.New(hold.Input{Session: fixture{test: test}.handle("session_123"), Owner: "one", Moment: moment, Window: 30 * time.Second})
	if failure != nil {
		test.Fatal(failure)
	}
	if held.Owner() != "one" || !held.Deadline().Equal(moment.Add(30*time.Second)) {
		test.Fatalf("hold = %+v", held)
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	if _, failure := hold.New(hold.Input{Session: fixture{test: test}.handle("session_123"), Owner: "one", Moment: time.Unix(1000, 0)}); failure == nil {
		test.Fatal("a hold without a window was accepted")
	}
}
