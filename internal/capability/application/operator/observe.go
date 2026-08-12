package operator

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
)

// passed is the stable, low-cardinality outcome label for a capability that completed without a refusal; a failure
// carries its own capability code instead.
const passed = "CAPABILITY_PASSED"

// monitor is the operator's port to operational observability. Its signatures are standard-library types, so the
// application core stays free of any telemetry vendor; an infrastructure adapter maps a scope to a trace span and a
// latency metric. A capability opens a scope through Watch, then reports its outcome once through the returned close.
type monitor interface {
	Watch(scope context.Context, capability string) (context.Context, func(string))
}

// watch instruments one capability execution: it carries the open telemetry scope, the capability name, and the start
// time. It holds only the capability name and a stable outcome label — never a device value — so no screenshot, element
// tree, typed text, or device identifier can reach a trace, metric, or log through this path.
type watch struct {
	operator   Operator
	close      func(string)
	capability string
	started    time.Time
}

// begin opens a telemetry scope for a capability and starts its latency measurement.
func (operator Operator) begin(scope context.Context, capability string) (context.Context, watch) {
	scope, close := operator.monitor.Watch(scope, capability)
	return scope, watch{operator: operator, close: close, capability: capability, started: time.Now()}
}

// finish closes the telemetry scope and logs the capability boundary. It is called through a deferred closure so it
// runs on every return path, including a refusal.
func (watch watch) finish(scope context.Context, failure error) {
	result := watch.result(failure)
	watch.close(result)
	watch.operator.logger.InfoContext(scope, "capability.completed",
		slog.String("capability", watch.capability),
		slog.String("outcome", result),
		slog.Duration("duration", time.Since(watch.started)),
	)
}

// result reduces the outcome to a stable, low-cardinality label — the capability's own failure code, or passed — never
// a device value.
func (watch watch) result(failure error) string {
	if failure == nil {
		return passed
	}
	if refusal, found := errors.AsType[Refusal](failure); found {
		return refusal.Code.String()
	}
	return outcome.Failed.String()
}
