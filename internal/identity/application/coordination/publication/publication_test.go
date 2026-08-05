package publication_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
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

func TestPublication(test *testing.T) {
	test.Parallel()

	notice, failure := publication.New(publication.Input{
		Session:  fixture{test: test}.handle("session_123"),
		Expected: epoch.Epoch(0),
		Key:      "key-1",
		Revision: 1,
		Moment:   time.Unix(1000, 0),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if notice.Key() != "key-1" || notice.Revision() != 1 {
		test.Fatalf("publication = %+v", notice)
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	if _, failure := publication.New(publication.Input{Session: fixture{test: test}.handle("session_123"), Key: "key-1", Moment: time.Unix(1000, 0)}); failure == nil {
		test.Fatal("a publication without a revision was accepted")
	}
}
