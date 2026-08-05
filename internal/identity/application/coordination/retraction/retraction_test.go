package retraction_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/retraction"
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

func TestRetraction(test *testing.T) {
	test.Parallel()

	request, failure := retraction.New(retraction.Input{
		Session: fixture{test: test}.handle("session_123"),
		Moment:  time.Unix(1000, 0),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if request.Session().String() != "session_123" {
		test.Fatalf("retraction = %+v", request)
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	if _, failure := retraction.New(retraction.Input{Session: fixture{test: test}.handle("session_123")}); failure == nil {
		test.Fatal("a retraction without a moment was accepted")
	}
}
