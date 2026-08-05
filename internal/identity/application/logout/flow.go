package logout

import (
	"context"
	"errors"

	fault "github.com/DrizzDev/platform/internal/identity/application/failure"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

// Flow clears the local session, then attempts bounded server revocation. The
// local clear is atomic and idempotent; a revocation failure is reported as a
// partial logout, never a lost local sign-out.
type Flow struct {
	vault       Vault
	publication Publication
	revocation  Revocation
	clock       Clock
}

func (flow Flow) Run(scope context.Context, _ Input) Result {
	record, failure := flow.vault.Active(scope)
	var missing Missing
	switch {
	case errors.As(failure, &missing):
		return Result{}
	case failure != nil:
		return flow.deny(failure)
	}
	if failure := flow.publication.Retract(scope, flow.clock.Now()); failure != nil {
		return flow.deny(failure)
	}
	if failure := flow.revocation.Revoke(scope, record); failure != nil {
		return flow.settle(code.Partial)
	}
	return Result{}
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
	return flow.settle(kind)
}

func (flow Flow) settle(kind code.Code) Result {
	value, _ := fault.New(fault.Input{Code: kind})
	return Result{failure: &value}
}
