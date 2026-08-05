package publication

import (
	"context"

	notice "github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/result"
	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

var _ login.Publication = Publisher{}

// Publisher implements the login Publication port. It resolves the session
// revision, writes the immutable vault candidate, and publishes it under a
// fenced compare-and-swap so concurrent logins cannot switch accounts.
type Publisher struct {
	vault   Vault
	ledger  Ledger
	random  login.Random
	session session.Session
}

func (publisher Publisher) Publish(scope context.Context, candidate login.Candidate) (login.Receipt, error) {
	if failure := publisher.admit(scope); failure != nil {
		return login.Receipt{}, failure
	}
	revision, failure := publisher.revision(scope, candidate)
	if failure != nil {
		return login.Receipt{}, failure
	}
	current, failure := publisher.ledger.Epoch(scope)
	if failure != nil {
		return login.Receipt{}, failure
	}
	handle, failure := publisher.mint()
	if failure != nil {
		return login.Receipt{}, failure
	}
	record, failure := publisher.record(draft{candidate: candidate, revision: revision, handle: handle})
	if failure != nil {
		return login.Receipt{}, login.Rejected{}
	}
	if failure := publisher.vault.Write(scope, record); failure != nil {
		return login.Receipt{}, login.Storage{}
	}
	request, failure := notice.New(notice.Input{
		Session:  publisher.session,
		Expected: current,
		Key:      record.Key().String(),
		Revision: revision,
		Moment:   candidate.Moment,
	})
	if failure != nil {
		return login.Receipt{}, login.Rejected{}
	}
	outcome, failure := publisher.ledger.Publish(scope, request)
	if failure != nil {
		publisher.discard(scope, record.Key())
		return login.Receipt{}, failure
	}
	switch outcome {
	case result.Published:
		return publisher.receipt(record), nil
	case result.Rejected:
		return publisher.contest(scope, candidate)
	case result.Uncertain:
		return login.Receipt{}, login.Unavailable{}
	}
	return login.Receipt{}, login.Unavailable{}
}

// discard best-effort removes the candidate a failed publish left in the vault.
// It ignores cancellation so a cancelled login still cleans up; a residual is
// reconciled at startup.
func (publisher Publisher) discard(scope context.Context, key credential.Key) {
	_ = publisher.vault.Delete(context.WithoutCancel(scope), key)
}

// admit refuses a new credential once the cleanup backlog is saturated, so a
// vault that keeps rejecting deletions cannot accumulate orphaned secrets
// without bound.
func (publisher Publisher) admit(scope context.Context) error {
	count, failure := publisher.ledger.Backlog(scope)
	if failure != nil {
		return login.Storage{}
	}
	if count >= saturation {
		return login.Storage{}
	}
	return nil
}
