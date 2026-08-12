package bridge

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// gauge instruments one round trip to the helper at the provider boundary: a span, a latency measurement, and a stable
// outcome. It carries only the wire method name and the outcome — never a device value or request payload — so nothing
// a device returns can reach a trace or metric through this path.
type gauge struct {
	channel *Channel
	span    trace.Span
	method  string
	started time.Time
}

// begin opens a span for one wire call and starts its latency measurement.
func (channel *Channel) begin(scope context.Context, method string) (context.Context, gauge) {
	scope, span := channel.tracer.Start(scope, method)
	return scope, gauge{channel: channel, span: span, method: method, started: time.Now()}
}

// close ends the span and records the latency against the method and its outcome. It runs through a deferred closure so
// every return path — including a timeout or an unavailable helper — is measured.
func (gauge gauge) close(scope context.Context, failure error) {
	result := gauge.grade(failure)
	gauge.span.SetAttributes(attribute.String("drizz.outcome", result))
	gauge.span.End()

	gauge.channel.duration.Record(scope, time.Since(gauge.started).Seconds(), metric.WithAttributes(
		attribute.String("drizz.method", gauge.method),
		attribute.String("drizz.outcome", result),
	))
}

// grade reduces a transport failure to a stable, low-cardinality outcome code, so the metric never carries an
// unbounded error string.
func (gauge gauge) grade(failure error) string {
	switch {
	case failure == nil:
		return "ok"
	case errors.Is(failure, context.DeadlineExceeded):
		return "timeout"
	case errors.Is(failure, context.Canceled):
		return "cancelled"
	case errors.Is(failure, Closed{}):
		return "closed"
	case errors.Is(failure, Unavailable{}):
		return "unavailable"
	default:
		return "failed"
	}
}
