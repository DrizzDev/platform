package sqlite

import (
	"database/sql"
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Store is the single, shared owner of the identity database. It is opened once
// at composition and injected wherever coordination state is read or written.
type Store struct {
	handle  *sql.DB
	logger  *slog.Logger
	tracer  trace.Tracer
	latency metric.Float64Histogram
}

func (store Store) Close() error {
	return store.handle.Close()
}
