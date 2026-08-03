package mcp

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type execution struct {
	server  Server
	started time.Time
	span    trace.Span
}

func (execution execution) finish(scope context.Context, failure error) error {
	result := execution.outcome(failure)
	if result == interrupted {
		execution.span.SetStatus(codes.Error, "MCP server stopped")
	}
	execution.span.SetAttributes(attribute.String("drizz.outcome", string(result)))
	execution.span.End()

	execution.record(scope, result)
	if result == success || result == cancelled {
		return nil
	}
	return failure
}
