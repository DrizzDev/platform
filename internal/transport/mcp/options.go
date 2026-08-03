package mcp

import (
	"io"
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/DrizzDev/platform/internal/application/release"
)

type Options struct {
	Release  release.Identity
	Logger   *slog.Logger
	External *slog.Logger
	Tracer   trace.Tracer
	Meter    metric.Meter
	Input    io.ReadCloser
	Output   io.Writer
}
