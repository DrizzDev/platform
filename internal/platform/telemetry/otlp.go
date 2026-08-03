package telemetry

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	metering "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracing "go.opentelemetry.io/otel/sdk/trace"
)

func (options Options) otlp(scope context.Context) (Provider, error) {
	traces, failure := otlptracehttp.New(
		scope,
		otlptracehttp.WithEndpointURL(options.Settings.Endpoint()),
	)
	if failure != nil {
		return Provider{}, failure
	}
	metrics, failure := otlpmetrichttp.New(
		scope,
		otlpmetrichttp.WithEndpointURL(options.Settings.Endpoint()),
	)
	if failure != nil {
		return Provider{}, errors.Join(failure, traces.Shutdown(scope))
	}
	identity := resource.NewSchemaless(
		attribute.String("service.name", options.Build.Name()),
		attribute.String("service.version", options.Build.Version()),
	)
	trace := tracing.NewTracerProvider(tracing.WithBatcher(traces), tracing.WithResource(identity))
	meter := metering.NewMeterProvider(
		metering.WithReader(metering.NewPeriodicReader(metrics)),
		metering.WithResource(identity),
	)
	return Provider{
		tracer:  trace.Tracer(options.Build.Name()),
		meter:   meter.Meter(options.Build.Name()),
		once:    &sync.Once{},
		outcome: new(error),
		close: func(scope context.Context) error {
			return errors.Join(meter.Shutdown(scope), trace.Shutdown(scope))
		},
	}, nil
}
