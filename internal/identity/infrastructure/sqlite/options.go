package sqlite

import (
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Options struct {
	Path   string
	Logger *slog.Logger
	Tracer trace.Tracer
	Meter  metric.Meter
}
