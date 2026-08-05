package observability

import (
	"context"
	"errors"
	"log/slog"

	"github.com/DrizzDev/platform/internal/platform/reporting"
	"github.com/DrizzDev/platform/internal/platform/telemetry"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Provider struct {
	diagnostics *slog.Logger
	sink        *slog.Logger
	external    *slog.Logger
	closer      func(context.Context) error
	reporting   reporting.Provider
	telemetry   telemetry.Provider
}

// Report records a runtime construction fault as its approved code alone. It
// carries no cause, so no caller content reaches diagnostics or the sink; the
// caller retains the underlying error for control flow.
func (provider Provider) Report(scope context.Context) {
	record := defect{name: command}
	provider.diagnostics.ErrorContext(scope, record.event())
	provider.sink.ErrorContext(scope, record.event())
}

func (provider Provider) Diagnostics() *slog.Logger {
	return provider.diagnostics
}

func (provider Provider) External() *slog.Logger {
	return provider.external
}

func (provider Provider) Tracer() trace.Tracer {
	return provider.telemetry.Tracer()
}

func (provider Provider) Meter() metric.Meter {
	return provider.telemetry.Meter()
}

func (provider Provider) Tracing() trace.TracerProvider {
	return provider.telemetry.Tracing()
}

func (provider Provider) Metering() metric.MeterProvider {
	return provider.telemetry.Metering()
}

// Close shuts observability down. A failure is recorded as a code alone; the
// underlying provider error is returned for control flow but never logged, so
// endpoint or machine detail cannot reach diagnostics.
func (provider Provider) Close(scope context.Context) error {
	failure := provider.terminate(scope)
	if failure != nil {
		provider.diagnostics.WarnContext(scope, "shutdown.failed")
	}
	return failure
}

func (provider Provider) terminate(scope context.Context) error {
	if provider.closer != nil {
		return provider.closer(scope)
	}
	return errors.Join(provider.reporting.Close(scope), provider.telemetry.Close(scope))
}
