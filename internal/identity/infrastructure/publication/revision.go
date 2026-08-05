package publication

import (
	"context"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

// revision returns the next revision for the session, rejecting a login for a
// different account than the one already active.
func (publisher Publisher) revision(scope context.Context, candidate login.Candidate) (uint64, error) {
	head, failure := publisher.ledger.Head(scope, publisher.session)
	var absent sqlite.Absent
	switch {
	case errors.As(failure, &absent):
		return 1, nil
	case failure != nil:
		return 0, failure
	}
	active, failure := publisher.vault.Read(scope, credential.Key(head.Key()))
	if failure != nil {
		return 0, login.Storage{}
	}
	if active.Subject().String() != candidate.Token.Subject.String() {
		return 0, login.Conflict{}
	}
	return head.Revision() + 1, nil
}

// contest resolves the outcome after losing the compare-and-swap. The same
// account winning concurrently is a success; a different account is a conflict.
func (publisher Publisher) contest(scope context.Context, candidate login.Candidate) (login.Receipt, error) {
	head, failure := publisher.ledger.Head(scope, publisher.session)
	if failure != nil {
		return login.Receipt{}, login.Unavailable{}
	}
	active, failure := publisher.vault.Read(scope, credential.Key(head.Key()))
	if failure != nil {
		return login.Receipt{}, login.Storage{}
	}
	if active.Subject().String() != candidate.Token.Subject.String() {
		return login.Receipt{}, login.Conflict{}
	}
	return publisher.receipt(active), nil
}
