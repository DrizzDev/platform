package mcp

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (execution execution) record(scope context.Context, result outcome) {
	elapsed := time.Since(execution.started)
	labels := metric.WithAttributes(
		attribute.String("drizz.operation", operation),
		attribute.String("drizz.outcome", string(result)),
	)
	execution.server.duration.Record(scope, elapsed.Seconds(), labels)
	execution.server.logger.InfoContext(
		scope,
		"mcp.completed",
		slog.String("operation", operation),
		slog.String("outcome", string(result)),
		slog.Duration("duration", elapsed),
	)
}
