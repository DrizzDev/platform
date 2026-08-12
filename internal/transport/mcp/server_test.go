package mcp_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"
	metering "go.opentelemetry.io/otel/sdk/metric"
	tracing "go.opentelemetry.io/otel/sdk/trace"

	"github.com/DrizzDev/platform/internal/application/release"
	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/transport/mcp"
)

type performer struct {
	shot   operator.Shot
	roster operator.Roster
	fail   error
}

func (performer performer) Screenshot(context.Context, operator.Target) (operator.Shot, error) {
	return performer.shot, performer.fail
}

func (performer performer) Snapshot(context.Context, operator.Target) (operator.Snapshot, error) {
	return operator.Snapshot{}, performer.fail
}

func (performer performer) Hierarchy(context.Context, operator.Target) (operator.Tree, error) {
	return operator.Tree{}, performer.fail
}

func (performer performer) Dimensions(context.Context, operator.Target) (operator.Extent, error) {
	return operator.Extent{}, performer.fail
}

func (performer performer) Devices(context.Context) (operator.Roster, error) {
	return performer.roster, performer.fail
}

func (performer performer) Install(context.Context, operator.Package) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Launch(context.Context, operator.Application) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Terminate(context.Context, operator.Application) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Wipe(context.Context, operator.Application) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Installed(context.Context, operator.Target) (operator.Listing, error) {
	return operator.Listing{}, performer.fail
}

func (performer performer) Running(context.Context, operator.Target) (operator.Listing, error) {
	return operator.Listing{}, performer.fail
}

func (performer performer) Foreground(context.Context, operator.Target) (operator.Report, error) {
	return operator.Report{}, performer.fail
}

func (performer performer) Url(context.Context, operator.Target) (operator.Report, error) {
	return operator.Report{}, performer.fail
}

func (performer performer) Disk(context.Context, operator.Target) (operator.Measure, error) {
	return operator.Measure{}, performer.fail
}

func (performer performer) Images(context.Context) (operator.Images, error) {
	return operator.Images{}, performer.fail
}

func (performer performer) Boot(context.Context, operator.Image) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Pause(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Resume(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Tap(context.Context, operator.Contact) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Swipe(context.Context, operator.Drag) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Pinch(context.Context, operator.Squeeze) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Press(context.Context, operator.Key) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Type(context.Context, operator.Entry) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Clear(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Back(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Home(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Locate(context.Context, operator.Fix) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

type dialog struct {
	test    *testing.T
	perform mcp.Perform
}

// session stands up the MCP server with the fake operator over an in-memory link and returns a connected agent.
func (dialog dialog) session(scope context.Context) *protocol.ClientSession {
	discard := slog.New(slog.NewTextHandler(io.Discard, nil))
	identity, failure := release.New(release.Input{Name: "drizz", Version: "1.0.0", Revision: "revision_123"})
	if failure != nil {
		dialog.test.Fatal(failure)
	}
	server, failure := mcp.New(mcp.Options{
		Logger:   discard,
		External: discard,
		Tracer:   tracing.NewTracerProvider().Tracer("test"),
		Meter:    metering.NewMeterProvider().Meter("test"),
		Release:  identity,
		Input:    io.NopCloser(strings.NewReader("")),
		Output:   io.Discard,
		Perform:  dialog.perform,
	})
	if failure != nil {
		dialog.test.Fatal(failure)
	}
	client, transport := protocol.NewInMemoryTransports()
	done := make(chan error, 1)
	go func() { done <- mcp.Serve(scope, mcp.Request{Server: server, Transport: transport}) }()
	agent := protocol.NewClient(&protocol.Implementation{Name: "test", Version: "1.0.0"}, nil)
	session, failure := agent.Connect(scope, client, nil)
	if failure != nil {
		dialog.test.Fatal(failure)
	}
	dialog.test.Cleanup(func() {
		_ = session.Close()
		<-done
	})
	return session
}

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
		Perform:  performer{},
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

func TestScreenshotTool(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := dialog{
		test:    test,
		perform: performer{shot: operator.Shot{Image: []byte("png-bytes"), Format: "PNG"}},
	}.session(scope)

	result, failure := session.CallTool(scope, &protocol.CallToolParams{
		Name:      "TakeScreenshot",
		Arguments: map[string]any{"serial": "s-1"},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if result.IsError || len(result.Content) != 1 {
		test.Fatalf("result = %#v", result)
	}
	image, valid := result.Content[0].(*protocol.ImageContent)
	if !valid {
		test.Fatalf("content = %#v", result.Content[0])
	}
	if string(image.Data) != "png-bytes" || image.MIMEType != "image/png" {
		test.Fatalf("image = %#v", image)
	}
}

func TestDevicesTool(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := dialog{
		test:    test,
		perform: performer{roster: operator.Roster{Serials: []string{"s-1", "s-2"}}},
	}.session(scope)

	result, failure := session.CallTool(scope, &protocol.CallToolParams{
		Name:      "ListDevices",
		Arguments: map[string]any{},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if result.IsError {
		test.Fatalf("result = %#v", result)
	}
	reported := fmt.Sprintf("%v", result.StructuredContent)
	if !strings.Contains(reported, "s-1") || !strings.Contains(reported, "s-2") {
		test.Fatalf("structured content = %q", reported)
	}
}

func TestScreenshotToolRefused(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := dialog{
		test:    test,
		perform: performer{fail: operator.Refusal{Code: outcome.Missing}},
	}.session(scope)

	result, failure := session.CallTool(scope, &protocol.CallToolParams{
		Name:      "TakeScreenshot",
		Arguments: map[string]any{"serial": "s-9"},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if !result.IsError || len(result.Content) != 1 {
		test.Fatalf("refusal = %#v", result)
	}
	text, valid := result.Content[0].(*protocol.TextContent)
	if !valid || !strings.Contains(text.Text, "not found") {
		test.Fatalf("content = %#v", result.Content[0])
	}
}

// TestCatalogCovered is the other half of the completeness gate: every catalogued capability must be offered as an
// agent tool, so a capability can never be catalogued but left off the agent connection.
func TestCatalogCovered(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := dialog{test: test, perform: performer{}}.session(scope)
	listed, failure := session.ListTools(scope, nil)
	if failure != nil {
		test.Fatal(failure)
	}
	present := make(map[string]bool, len(listed.Tools))
	for _, tool := range listed.Tools {
		present[tool.Name] = true
	}
	for _, entry := range catalog.New().List() {
		if !present[entry.Title()] {
			test.Errorf("catalogued capability %q has no agent tool", entry.Title())
		}
	}
}

func TestDependencies(test *testing.T) {
	test.Parallel()

	_, failure := mcp.New(mcp.Options{})
	if failure == nil {
		test.Fatal("missing dependencies were accepted")
	}
}
