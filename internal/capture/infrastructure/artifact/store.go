package artifact

import (
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Store is the single owner of the content-addressed artifact directory,
// opened once at composition and injected wherever artifacts are written or read.
type Store struct {
	ceiling int64
	root    string
	temp    string
	logger  *slog.Logger
	tracer  trace.Tracer
	latency metric.Float64Histogram
}
