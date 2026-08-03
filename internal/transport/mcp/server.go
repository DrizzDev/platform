package mcp

import (
	"context"
	"io"
	"log/slog"
	"time"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	server   *protocol.Server
	logger   *slog.Logger
	tracer   trace.Tracer
	duration metric.Float64Histogram
	input    io.ReadCloser
	output   io.Writer
}

func (server Server) Run(scope context.Context) error {
	return server.serve(scope, transport{input: server.input, output: server.output})
}

func (server Server) serve(scope context.Context, transport protocol.Transport) error {
	scope, span := server.tracer.Start(scope, operation)
	run := execution{server: server, span: span, started: time.Now()}
	server.logger.InfoContext(scope, "mcp.started", slog.String("operation", operation))
	failure := server.server.Run(scope, transport)
	return run.finish(scope, failure)
}
