package login

import (
	"context"
	"errors"

	fault "github.com/DrizzDev/platform/internal/identity/application/failure"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
	"github.com/DrizzDev/platform/internal/identity/domain/standing"
)

// Flow completes a sign-in: it establishes a validated provider token through the
// chosen front-end, then publishes the credential under a fenced compare-and-swap.
type Flow struct {
	establishment Establishment
	publication   Publication
	authority     Authority
	clock         Clock
}

func (flow Flow) Run(scope context.Context, _ Input) Result {
	token, failure := flow.establishment.Establish(scope)
	if failure != nil {
		return flow.deny(failure)
	}
	receipt, failure := flow.publication.Publish(scope, Candidate{Token: token, Moment: flow.clock.Now()})
	if failure != nil {
		return flow.deny(failure)
	}
	tenant, failure := flow.authority.Authorize(scope, Grant{Token: token.Access, Expiry: token.Expiry})
	var forbidden Forbidden
	if errors.As(failure, &forbidden) {
		_ = flow.publication.Retract(scope, flow.clock.Now())
		return flow.deny(failure)
	}
	// A non-forbidden authority failure is best-effort: local sign-in proceeds
	// and the cloud re-authorizes every later operation.
	return Result{
		subject:  receipt.Subject,
		session:  receipt.Session,
		method:   receipt.Method,
		standing: standing.Active,
		expiry:   receipt.Expiry,
		tenant:   tenant,
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
