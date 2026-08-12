// Package telemetry adapts the integration application's observability port to OpenTelemetry: it opens a span per
// operation and records the operation's latency, keeping the application core free of any telemetry vendor.
package telemetry

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Monitor maps the integration application's telemetry scope to a trace span and a latency metric. Its labels are the
// operation name and a stable outcome label only — never an agent identifier, path, or a person's prompt — so telemetry
// stays low-cardinality and carries no captured content.
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
	duration, failure := options.Meter.Float64Histogram("drizz.integration.duration", metric.WithUnit("s"))
	if failure != nil {
		return Monitor{}, failure
	}
	return Monitor{tracer: options.Tracer, duration: duration}, nil
}

// Watch opens a span for an operation and returns a close that ends the span and records its latency against the
// operation and its outcome. The close is called once, on the operation's return path.
func (monitor Monitor) Watch(scope context.Context, operation string) (context.Context, func(string)) {
	scope, span := monitor.tracer.Start(scope, operation)
	started := time.Now()
	return scope, func(outcome string) {
		span.SetAttributes(attribute.String("drizz.outcome", outcome))
		span.End()
		monitor.duration.Record(scope, time.Since(started).Seconds(), metric.WithAttributes(
			attribute.String("drizz.operation", operation),
			attribute.String("drizz.outcome", outcome),
		))
	}
}
