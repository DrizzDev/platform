package session

import (
	"context"
	"errors"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/attempt"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/marking"
	fault "github.com/DrizzDev/platform/internal/identity/application/failure"
	"github.com/DrizzDev/platform/internal/identity/application/grant"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
	"github.com/DrizzDev/platform/internal/identity/domain/standing"
)

// Flow renews the active session in place: it fences the active revision once,
// exchanges the refresh token, and publishes the rotated credential under a
// fenced compare-and-swap. A missing, fenced, or unrenewable credential resolves
// into a required sign-in rather than a silent retry.
type Flow struct {
	vault       Vault
	refresh     Refresh
	publication Publication
	epoch       Epoch
	clock       Clock
}

// yield is the internal result of one renewal: the trusted context and the
// rotated secrets it was drawn from.
type yield struct {
	receipt Receipt
	renewal Renewal
}

// Current renews the session and returns its trusted context.
func (flow Flow) Current(scope context.Context, _ Input) Result {
	outcome, failure := flow.renew(scope)
	if failure != nil {
		return flow.deny(failure)
	}
	return flow.settle(outcome.receipt)
}

// Access renews the session and returns the access token for a cloud call. The
// grant is memory-only and confined; it is never persisted or returned to a
// transport.
func (flow Flow) Access(scope context.Context, _ Input) (grant.Credential, error) {
	outcome, failure := flow.renew(scope)
	if failure != nil {
		return grant.Credential{}, failure
	}
	return grant.New(grant.Input{Token: outcome.renewal.Access, Expiry: outcome.renewal.Expiry})
}

func (flow Flow) renew(scope context.Context) (yield, error) {
	record, failure := flow.vault.Active(scope)
	if failure != nil {
		return yield{}, failure
	}
	current, failure := flow.epoch.Read(scope)
	if failure != nil {
		return yield{}, failure
	}
	entry, failure := attempt.New(attempt.Input{Revision: record.Revision(), Epoch: current})
	if failure != nil {
		return yield{}, failure
	}
	mark, failure := marking.New(marking.Input{Session: record.Session(), Attempt: entry})
	if failure != nil {
		return yield{}, failure
	}
	if failure := flow.publication.Attempt(scope, mark); failure != nil {
		return yield{}, failure
	}
	renewal, failure := flow.refresh.Renew(scope, record)
	if failure != nil {
		return yield{}, failure
	}
	receipt, failure := flow.publication.Publish(scope, Candidate{
		Prior:    record,
		Renewal:  renewal,
		Expected: current,
		Moment:   flow.clock.Now(),
	})
	if failure != nil {
		return yield{}, failure
	}
	return yield{receipt: receipt, renewal: renewal}, nil
}

func (flow Flow) settle(receipt Receipt) Result {
	return Result{
		subject:  receipt.Subject,
		session:  receipt.Session,
		method:   receipt.Method,
		standing: standing.Active,
		expiry:   receipt.Expiry,
	}
}

func (flow Flow) deny(cause error) Result {
	kind := code.Failed
	var carrier interface{ reason() code.Code }
	switch {
	case errors.Is(cause, context.Canceled), errors.Is(cause, context.DeadlineExceeded):
		kind = code.Cancelled
	case errors.As(cause, &carrier):
		kind = carrier.reason()
	}
	value, _ := fault.New(fault.Input{Code: kind})
	return Result{failure: &value}
}
