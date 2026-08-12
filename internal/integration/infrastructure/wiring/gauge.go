package wiring

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
)

// gauge instruments one configuration-file operation at the filesystem provider boundary: a span, a latency
// measurement, and a stable outcome. It carries only the operation name and the outcome — never a file path, an agent
// identifier, or file contents — so nothing an agent's configuration holds can reach a trace or metric through here.
type gauge struct {
	store     Store
	span      trace.Span
	operation string
	started   time.Time
}

// begin opens a span for one store operation and starts its latency measurement.
func (store Store) begin(scope context.Context, operation string) (context.Context, gauge) {
	scope, span := store.tracer.Start(scope, operation)
	return scope, gauge{store: store, span: span, operation: operation, started: time.Now()}
}

// close ends the span and records the latency against the operation and its outcome. It runs through a deferred closure
// so every return path — a missing file, a malformed file, a locked file — is measured.
func (gauge gauge) close(scope context.Context, failure error) {
	result := gauge.grade(failure)
	gauge.span.SetAttributes(attribute.String("drizz.outcome", result))
	gauge.span.End()
	gauge.store.duration.Record(scope, time.Since(gauge.started).Seconds(), metric.WithAttributes(
		attribute.String("drizz.operation", gauge.operation),
		attribute.String("drizz.outcome", result),
	))
}

// grade reduces a store failure to a stable, low-cardinality outcome code, so the metric never carries an unbounded
// error string.
func (gauge gauge) grade(failure error) string {
	if failure == nil {
		return "ok"
	}
	if _, found := errors.AsType[connect.Malformed](failure); found {
		return "malformed"
	}
	if _, found := errors.AsType[connect.Unsupported](failure); found {
		return "unsupported"
	}
	if _, found := errors.AsType[connect.Occupied](failure); found {
		return "occupied"
	}
	if _, found := errors.AsType[connect.Locked](failure); found {
		return "locked"
	}
	return "failed"
}
