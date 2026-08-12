package connect

import (
	"context"
	"log/slog"
	"time"
)

// monitor is the installer's port to operational observability. Its signatures are standard-library types, so the
// application core stays free of any telemetry vendor; an infrastructure adapter maps a scope to a trace span and a
// latency metric. An operation opens a scope through Watch, then reports its outcome once through the returned close.
type monitor interface {
	Watch(scope context.Context, operation string) (context.Context, func(string))
}

// watch instruments one installer operation: the open telemetry scope, the operation name, and the start time. It
// holds only the operation name and a stable outcome label — never an agent identifier or a person's prompt — so no
// captured content can reach a trace, metric, or log through this path.
type watch struct {
	installer Installer
	close     func(string)
	operation string
	started   time.Time
}

// begin opens a telemetry scope for an operation and starts its latency measurement.
func (installer Installer) begin(scope context.Context, operation string) (context.Context, watch) {
	scope, close := installer.monitor.Watch(scope, operation)
	return scope, watch{installer: installer, close: close, operation: operation, started: time.Now()}
}

// finish closes the telemetry scope and logs the operation boundary. It is called through a deferred closure so it runs
// on every return path.
func (watch watch) finish(scope context.Context, result string) {
	watch.close(result)
	watch.installer.logger.InfoContext(scope, "integration.completed",
		slog.String("operation", watch.operation),
		slog.String("outcome", result),
		slog.Duration("duration", time.Since(watch.started)),
	)
}

// summary is one operation's result to reduce to an outcome label: the per-agent report and the top-level failure, if
// the operation could not start at all.
type summary struct {
	report  Report
	failure error
}

// tally reduces a report and a top-level failure to one stable, low-cardinality operation outcome: failed when the
// operation could not start, partial when some agents failed, done otherwise.
func (Installer) tally(result summary) string {
	if result.failure != nil {
		return "failed"
	}
	for _, outcome := range result.report.outcomes {
		if outcome.state == Failed {
			return "partial"
		}
	}
	return "done"
}
