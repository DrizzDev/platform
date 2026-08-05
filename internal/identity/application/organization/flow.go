package organization

import (
	"context"
	"errors"

	fault "github.com/DrizzDev/platform/internal/identity/application/failure"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

// Flow resolves the current organization for the authenticated subject through
// Drizz Cloud. The cloud is the sole authority; a local claim or request never
// substitutes for its decision.
type Flow struct {
	resolver Resolver
}

func (flow Flow) Resolve(scope context.Context, _ Input) Result {
	tenant, failure := flow.resolver.Resolve(scope)
	if failure != nil {
		return flow.deny(failure)
	}
	return Result{organization: tenant}
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
