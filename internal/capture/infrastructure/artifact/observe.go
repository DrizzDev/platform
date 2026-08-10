package artifact

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

const (
	outcome    = "drizz.outcome"
	operation  = "drizz.operation"
	boundary   = "capture.artifact"
	instrument = "drizz.capture.duration"
)

type verdict string

const (
	settled verdict = "OK"
	faulted verdict = "FAULT"
)

type reader = func(context.Context) error

type probe struct {
	work reader
	name string
}

// observe wraps a file operation with a span and a duration metric. A fault records only the operation code, never the
// filesystem cause, which stays with the returned error for the application boundary to map.
func (store Store) observe(scope context.Context, probe probe) (failure error) {
	scope, span := store.tracer.Start(scope, boundary+"."+probe.name)
	start := time.Now()

	defer func() {
		result := store.grade(failure)
		if result == faulted {
			span.SetStatus(codes.Error, probe.name)
			store.logger.ErrorContext(scope, "capture.artifact.failed", slog.String(operation, probe.name))
		}
		store.latency.Record(scope, time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String(operation, probe.name),
			attribute.String(outcome, string(result))))
		span.End()
	}()

	failure = probe.work(scope)
	return failure
}

func (Store) grade(failure error) verdict {
	if failure == nil {
		return settled
	}
	return faulted
}
