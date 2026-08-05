package publication_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/attempt"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	"github.com/DrizzDev/platform/internal/identity/application/session"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/publication"
)

func TestRenew(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := fixture.store()
	safe := fixture.safe(&locker{entries: map[string][]byte{}})
	handle := fixture.handle()
	publisher := fixture.publisher(publication.Options{Vault: safe, Ledger: store, Random: entropy{}, Session: handle})
	scope := context.Background()

	if _, failure := publisher.Publish(scope, fixture.candidate("google-oauth2|first")); failure != nil {
		test.Fatal(failure)
	}

	renewer := publisher.Renewer()
	prior, failure := renewer.Active(scope)
	if failure != nil {
		test.Fatal(failure)
	}
	entry, failure := attempt.New(attempt.Input{Revision: prior.Revision(), Epoch: fixture.current(scope, renewer)})
	if failure != nil {
		test.Fatal(failure)
	}
	mark, failure := marking.New(marking.Input{Session: handle, Attempt: entry})
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := renewer.Attempt(scope, mark); failure != nil {
		test.Fatal(failure)
	}

	receipt, failure := renewer.Publish(scope, session.Candidate{
		Prior:    prior,
		Renewal:  session.Renewal{Refresh: []byte("rotated-refresh"), Expiry: time.Unix(4000, 0)},
		Expected: fixture.current(scope, renewer),
		Moment:   time.Unix(3000, 0),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if receipt.Subject.String() != "google-oauth2|first" {
		test.Fatalf("receipt = %+v", receipt)
	}

	head, failure := store.Head(scope, handle)
	if failure != nil {
		test.Fatal(failure)
	}
	if head.Revision() != 2 || !strings.HasPrefix(head.Key(), "LOCAL#2#") {
		test.Fatalf("head = %+v", head)
	}
}

func (fixture fixture) current(scope context.Context, renewer publication.Refresher) epoch.Epoch {
	fixture.test.Helper()
	value, failure := renewer.Read(scope)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return value
}
