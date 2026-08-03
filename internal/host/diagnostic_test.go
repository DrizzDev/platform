package host_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/host"
)

func TestTransport(test *testing.T) {
	test.Parallel()

	var diagnostics bytes.Buffer
	process, failure := host.New(host.Options{
		Arguments: []string{"mcp"},
		Streams: host.Streams{
			Input:   faulty{},
			Output:  io.Discard,
			Failure: &diagnostics,
		},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := process.Run(context.Background()); failure == nil {
		test.Fatal("a transport failure was not surfaced")
	}
	report := diagnostics.String()
	if strings.Contains(report, `"level":"ERROR"`) || strings.Contains(report, "command.failed") {
		test.Fatalf("an external transport failure was reported as an internal defect: %q", report)
	}
	if !strings.Contains(report, `"outcome":"INTERRUPTED"`) {
		test.Fatalf("transport failure was not recorded: %q", report)
	}
}

func TestUsage(test *testing.T) {
	test.Parallel()

	var output, diagnostics bytes.Buffer
	process, failure := host.New(host.Options{
		Arguments: []string{"bogus"},
		Streams: host.Streams{
			Input:   io.NopCloser(strings.NewReader("")),
			Output:  &output,
			Failure: &diagnostics,
		},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := process.Run(context.Background()); failure == nil {
		test.Fatal("a usage mistake was not surfaced")
	}
	report := diagnostics.String()
	if !strings.Contains(report, "unknown command") {
		test.Fatalf("usage mistake was not shown to the user: %q", report)
	}
	if strings.Contains(report, "command.failed") || strings.Contains(report, `"level":"ERROR"`) {
		test.Fatalf("usage mistake was escalated to a reported failure: %q", report)
	}
}

func TestCancellation(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithCancel(context.Background())
	cancel()

	var diagnostics bytes.Buffer
	process, failure := host.New(host.Options{
		Arguments: []string{"mcp"},
		Streams: host.Streams{
			Input:   io.NopCloser(strings.NewReader("")),
			Output:  io.Discard,
			Failure: &diagnostics,
		},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := process.Run(scope); failure != nil {
		test.Fatalf("cancellation was not a clean shutdown: %v", failure)
	}
	report := diagnostics.String()
	if strings.Contains(report, "command.failed") || strings.Contains(report, `"level":"ERROR"`) {
		test.Fatalf("cancellation was reported as a command failure: %q", report)
	}
	if !strings.Contains(report, `"message":"mcp.completed"`) {
		test.Fatalf("clean shutdown was not recorded: %q", report)
	}
}

type faulty struct{}

func (faulty) Read([]byte) (int, error) {
	return 0, errors.New("device read failure")
}

func (faulty) Close() error {
	return nil
}
