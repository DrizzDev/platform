package artifact

import (
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Options struct {
	Ceiling int64
	Root    string
	Logger  *slog.Logger
	Tracer  trace.Tracer
	Meter   metric.Meter
}
