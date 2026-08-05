package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Provider struct {
	tracer  trace.Tracer
	meter   metric.Meter
	traces  trace.TracerProvider
	meters  metric.MeterProvider
	close   func(context.Context) error
	once    *sync.Once
	outcome *error
}

func (provider Provider) Tracer() trace.Tracer {
	return provider.tracer
}

func (provider Provider) Meter() metric.Meter {
	return provider.meter
}

func (provider Provider) Tracing() trace.TracerProvider {
	return provider.traces
}

func (provider Provider) Metering() metric.MeterProvider {
	return provider.meters
}

func (provider Provider) Close(scope context.Context) error {
	if provider.close == nil {
		return nil
	}
	provider.once.Do(func() {
		*provider.outcome = provider.close(scope)
	})
	return *provider.outcome
}
