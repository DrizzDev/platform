package telemetry

import (
	"context"
	"errors"
	"sync"

	metering "go.opentelemetry.io/otel/sdk/metric"
	tracing "go.opentelemetry.io/otel/sdk/trace"
)

func (options Options) noop() Provider {
	traces := tracing.NewTracerProvider()
	meters := metering.NewMeterProvider()
	return Provider{
		tracer:  traces.Tracer(options.Build.Name()),
		meter:   meters.Meter(options.Build.Name()),
		once:    &sync.Once{},
		outcome: new(error),
		close: func(scope context.Context) error {
			return errors.Join(meters.Shutdown(scope), traces.Shutdown(scope))
		},
	}
}
