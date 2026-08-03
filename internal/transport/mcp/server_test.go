package mcp_test

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"
	metering "go.opentelemetry.io/otel/sdk/metric"
	tracing "go.opentelemetry.io/otel/sdk/trace"

	"github.com/DrizzDev/platform/internal/application/release"
	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/transport/mcp"
)

func TestNegotiation(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, transport := protocol.NewInMemoryTransports()
	traces := tracing.NewTracerProvider()
	meters := metering.NewMeterProvider()
	info := build.Read()
	identity, failure := release.New(release.Input{
		Name:     info.Name(),
		Version:  info.Version(),
		Revision: info.Revision(),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	discard := slog.New(slog.NewTextHandler(io.Discard, nil))
	server, failure := mcp.New(mcp.Options{
		Logger:   discard,
		External: discard,
		Tracer:   traces.Tracer("test"),
		Meter:    meters.Meter("test"),
		Release:  identity,
		Input:    io.NopCloser(strings.NewReader("")),
		Output:   io.Discard,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	done := make(chan error, 1)
	go func() {
		done <- mcp.Serve(scope, mcp.Request{Server: server, Transport: transport})
	}()

	agent := protocol.NewClient(&protocol.Implementation{Name: "test", Version: "1.0.0"}, nil)
	session, failure := agent.Connect(scope, client, nil)
	if failure != nil {
		test.Fatal(failure)
	}
	test.Cleanup(func() {
		if failure := session.Close(); failure != nil {
			test.Error(failure)
		}
	})
	if failure := session.Ping(scope, nil); failure != nil {
		test.Fatal(failure)
	}
	implementation := session.InitializeResult().ServerInfo
	if implementation.Name != identity.Name() || implementation.Version != identity.Version() {
		test.Fatalf("server = %#v", implementation)
	}

	cancel()
	if failure := <-done; failure != nil {
		test.Fatalf("cancellation was not a clean shutdown: %v", failure)
	}
}

func TestDependencies(test *testing.T) {
	test.Parallel()

	_, failure := mcp.New(mcp.Options{})
	if failure == nil {
		test.Fatal("missing dependencies were accepted")
	}
}
