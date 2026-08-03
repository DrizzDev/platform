package host_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/host"
)

func TestVersion(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	process, failure := host.New(host.Options{
		Arguments:   []string{"version"},
		Environment: []string{"DRIZZ_UNKNOWN=value"},
		Streams: host.Streams{
			Input:   io.NopCloser(strings.NewReader("")),
			Output:  &output,
			Failure: io.Discard,
		},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := process.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if !strings.HasPrefix(output.String(), "drizz ") {
		test.Fatalf("output = %q", output.String())
	}
}

func TestConfiguration(test *testing.T) {
	test.Parallel()

	var diagnostics bytes.Buffer
	process, failure := host.New(host.Options{
		Arguments:   []string{"mcp"},
		Environment: []string{"DRIZZ_UNKNOWN=value"},
		Streams: host.Streams{
			Input:   io.NopCloser(strings.NewReader("")),
			Output:  io.Discard,
			Failure: &diagnostics,
		},
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := process.Run(context.Background()); failure == nil {
		test.Fatal("invalid configuration was accepted")
	}
	if !strings.Contains(diagnostics.String(), "unknown Drizz setting") {
		test.Fatalf("configuration error was not surfaced: %q", diagnostics.String())
	}
}
