package publication

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	notice "github.com/DrizzDev/platform/internal/identity/application/coordination/publication"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/result"
	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/application/session"
	handle "github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

var (
	_ session.Vault       = Refresher{}
	_ session.Epoch       = Refresher{}
	_ session.Publication = Refresher{}
)

// Refresher implements the session coordination ports. It reads the active
// credential, fences its revision one time, and publishes the rotated candidate
// under a fenced compare-and-swap so a concurrent login, logout, or renewal can
// never resurrect a superseded credential.
type Refresher struct {
	vault   Vault
	ledger  Ledger
	random  login.Random
	session handle.Session
}

// Renewer yields the session adapter over the same vault and coordination store
// as the publisher, so renewal shares one credential-coordination boundary.
func (publisher Publisher) Renewer() Refresher {
	return Refresher(publisher)
}

// Active reads the credential the active pointer names, reporting Missing when
// nothing is signed in.
func (refresher Refresher) Active(scope context.Context) (credential.Record, error) {
	head, failure := refresher.ledger.Head(scope, refresher.session)
	var absent sqlite.Absent
	switch {
	case errors.As(failure, &absent):
		return credential.Record{}, session.Missing{}
	case failure != nil:
		return credential.Record{}, session.Storage{}
	}
	record, failure := refresher.vault.Read(scope, credential.Key(head.Key()))
	if failure != nil {
		return credential.Record{}, session.Storage{}
	}
	return record, nil
}

func (refresher Refresher) Read(scope context.Context) (epoch.Epoch, error) {
	current, failure := refresher.ledger.Epoch(scope)
	if failure != nil {
		return 0, session.Storage{}
	}
	return current, nil
}

// Attempt marks the active revision one time. A replayed or superseded revision
// is uncertain: the refresh token behind it may already have been spent.
func (refresher Refresher) Attempt(scope context.Context, mark marking.Marking) error {
	failure := refresher.ledger.Fence(scope, mark)
	var fenced sqlite.Fenced
	var contention sqlite.Contention
	switch {
	case errors.As(failure, &fenced), errors.As(failure, &contention):
		return session.Uncertain{}
	case failure != nil:
		return session.Storage{}
	default:
		return nil
	}
}

// Publish writes the rotated candidate and advances the head under the fenced
// epoch. A lost compare-and-swap means the session changed underneath the
// renewal, so the credential is uncertain and a new sign-in is required.
func (refresher Refresher) Publish(scope context.Context, candidate session.Candidate) (session.Receipt, error) {
	if failure := refresher.admit(scope); failure != nil {
		return session.Receipt{}, failure
	}
	record, failure := refresher.assemble(candidate)
	if failure != nil {
		return session.Receipt{}, session.Uncertain{}
	}
	if failure := refresher.vault.Write(scope, record); failure != nil {
		return session.Receipt{}, session.Storage{}
	}
	request, failure := notice.New(notice.Input{
		Session:  refresher.session,
		Expected: candidate.Expected,
		Key:      record.Key().String(),
		Revision: record.Revision(),
		Moment:   candidate.Moment,
	})
	if failure != nil {
		refresher.discard(scope, record.Key())
		return session.Receipt{}, session.Uncertain{}
	}
	outcome, failure := refresher.ledger.Publish(scope, request)
	if failure != nil {
		refresher.discard(scope, record.Key())
		return session.Receipt{}, session.Uncertain{}
	}
	if outcome != result.Published {
		return session.Receipt{}, session.Uncertain{}
	}
	return refresher.receipt(record), nil
}

// assemble builds the rotated credential from the prior record and the renewed
// secret, drawing a fresh unique handle and advancing the revision.
func (refresher Refresher) assemble(candidate session.Candidate) (credential.Record, error) {
	prior := candidate.Prior
	name, failure := refresher.mint()
	if failure != nil {
		return credential.Record{}, failure
	}
	return credential.New(credential.Input{
		Issuer:   prior.Issuer(),
		Client:   prior.Client(),
		Handle:   name,
		Subject:  prior.Subject(),
		Session:  prior.Session(),
		Method:   prior.Method(),
		Refresh:  candidate.Renewal.Refresh,
		Issued:   candidate.Moment,
		Expiry:   candidate.Renewal.Expiry,
		Revision: prior.Revision() + 1,
		Schema:   prior.Schema(),
	})
}

func (refresher Refresher) receipt(record credential.Record) session.Receipt {
	return session.Receipt{
		Subject: record.Subject(),
		Session: record.Session(),
		Method:  record.Method(),
		Expiry:  record.Expiry(),
	}
}

// mint draws a unique per-renewal handle so a concurrent renewal writes a
// distinct vault candidate key and can never overwrite another attempt.
func (refresher Refresher) mint() (string, error) {
	bytes, failure := refresher.random.Bytes(entropy)
	if failure != nil {
		return "", failure
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// discard best-effort removes the candidate a failed publish left in the vault,
// ignoring cancellation so a cancelled renewal still cleans up.
func (refresher Refresher) discard(scope context.Context, key credential.Key) {
	_ = refresher.vault.Delete(context.WithoutCancel(scope), key)
}

// admit refuses a rotated credential once the cleanup backlog is saturated, so a
// vault that keeps rejecting deletions cannot accumulate orphaned secrets
// without bound.
func (refresher Refresher) admit(scope context.Context) error {
	count, failure := refresher.ledger.Backlog(scope)
	if failure != nil {
		return session.Storage{}
	}
	if count >= saturation {
		return session.Storage{}
	}
	return nil
}
