// Package telemetry adapts the operator's observability port to OpenTelemetry: it opens a span per capability and
// records the capability's latency, keeping the application core free of any telemetry vendor.
package telemetry

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Monitor maps the operator's telemetry scope to a trace span and a latency metric. Its labels are the capability name
// and the stable outcome label only — never a device serial, path, or payload — so telemetry stays low-cardinality and
// carries no device content.
type Monitor struct {
	tracer   trace.Tracer
	duration metric.Float64Histogram
}

type Options struct {
	Tracer trace.Tracer
	Meter  metric.Meter
}

func New(options Options) (Monitor, error) {
	if options.Tracer == nil || options.Meter == nil {
		return Monitor{}, errors.New("telemetry tracer and meter are required")
	}
	duration, failure := options.Meter.Float64Histogram("drizz.capability.duration", metric.WithUnit("s"))
	if failure != nil {
		return Monitor{}, failure
	}
	return Monitor{tracer: options.Tracer, duration: duration}, nil
}

// Watch opens a span for a capability and returns a close that ends the span and records its latency against the
// capability and its outcome. The close is called once, on the capability's return path.
func (monitor Monitor) Watch(scope context.Context, capability string) (context.Context, func(string)) {
	scope, span := monitor.tracer.Start(scope, capability)
	started := time.Now()
	return scope, func(outcome string) {
		span.SetAttributes(attribute.String("drizz.outcome", outcome))
		span.End()
		monitor.duration.Record(scope, time.Since(started).Seconds(), metric.WithAttributes(
			attribute.String("drizz.capability", capability),
			attribute.String("drizz.outcome", outcome),
		))
	}
}
