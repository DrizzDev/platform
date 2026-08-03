package observability_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/platform/configuration"
	"github.com/DrizzDev/platform/internal/platform/observability"
)

const canary = "canary_secret_value"

func TestReporting(test *testing.T) {
	test.Parallel()

	var diagnostics, sink bytes.Buffer
	provider := observability.Wire(observability.Wiring{
		Diagnostics: slog.New(slog.NewJSONHandler(&diagnostics, nil)),
		Sink:        slog.NewJSONHandler(&sink, nil),
	})
	provider.Report(context.Background())

	for name, output := range map[string]*bytes.Buffer{"diagnostics": &diagnostics, "sink": &sink} {
		var record map[string]any
		if failure := json.Unmarshal(output.Bytes(), &record); failure != nil {
			test.Fatalf("%s: %v", name, failure)
		}
		if record["level"] != "ERROR" || record["msg"] != "command.failed" {
			test.Fatalf("%s did not emit the code: %v", name, record)
		}
		if record["error"] != nil {
			test.Fatalf("%s carried a cause: %v", name, record["error"])
		}
	}
}

func TestPrivacy(test *testing.T) {
	test.Parallel()

	// The reporting API accepts no cause, so a canary cannot be passed. This
	// guards that a future signature change cannot reintroduce a content channel.
	var diagnostics, sink bytes.Buffer
	provider := observability.Wire(observability.Wiring{
		Diagnostics: slog.New(slog.NewJSONHandler(&diagnostics, nil)),
		Sink:        slog.NewJSONHandler(&sink, nil),
	})
	provider.Report(context.Background())

	for _, output := range []*bytes.Buffer{&diagnostics, &sink} {
		if bytes.Contains(output.Bytes(), []byte(canary)) {
			test.Fatalf("a report carried content: %q", output.String())
		}
	}
}

func TestShutdown(test *testing.T) {
	test.Parallel()

	var diagnostics bytes.Buffer
	provider := observability.Wire(observability.Wiring{
		Diagnostics: slog.New(slog.NewJSONHandler(&diagnostics, nil)),
		Sink:        slog.NewJSONHandler(&diagnostics, nil),
		Close: func(context.Context) error {
			return errors.New("otlp flush to https://collector.example failed: " + canary)
		},
	})

	failure := provider.Close(context.Background())
	if failure == nil {
		test.Fatal("the shutdown error was not returned for control flow")
	}
	if bytes.Contains(diagnostics.Bytes(), []byte(canary)) {
		test.Fatalf("a provider shutdown error reached diagnostics: %q", diagnostics.String())
	}
	if !bytes.Contains(diagnostics.Bytes(), []byte("shutdown.failed")) {
		test.Fatalf("the shutdown fault was not recorded as a code: %q", diagnostics.String())
	}
}

func TestIsolation(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	provider := fixture{test: test, output: &output}.provider(nil)
	provider.External().Error("server run cancelled", slog.String("session_id", ""))
	if failure := provider.Close(context.Background()); failure != nil {
		test.Error(failure)
	}

	if output.Len() != 0 {
		test.Fatalf("external SDK diagnostics reached Drizz output: %q", output.String())
	}
}

type fixture struct {
	test   *testing.T
	output *bytes.Buffer
}

func (fixture fixture) provider(environment []string) observability.Provider {
	fixture.test.Helper()
	settings, failure := configuration.New(environment).Load()
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	provider, failure := observability.New(context.Background(), observability.Options{
		Build:    build.Read(),
		Settings: settings,
		Output:   fixture.output,
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return provider
}
