package transport

import (
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Options configures a resilient, observable HTTP client. The same builder
// serves every outbound integration; a caller supplies its own timeouts, retry
// budget, body ceiling, and observability providers, and chooses whether trace
// context is propagated (only ever to first-party services).
type Options struct {
	Timeout   time.Duration
	Dial      time.Duration
	Ceiling   int64
	Retries   int
	Minimum   time.Duration
	Maximum   time.Duration
	Tracing   trace.TracerProvider
	Metering  metric.MeterProvider
	Propagate bool
}

// propagator injects W3C trace context only when the target is first-party. A
// third-party issuer receives no propagator, so internal trace identifiers never
// leave the process.
func (options Options) propagator() propagation.TextMapPropagator {
	if options.Propagate {
		return propagation.TraceContext{}
	}
	return propagation.NewCompositeTextMapPropagator()
}
