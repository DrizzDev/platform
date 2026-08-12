package host_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/host"
)

func TestScreenshotUnprepared(test *testing.T) {
	test.Parallel()

	var diagnostics bytes.Buffer
	process, failure := host.New(host.Options{
		Arguments: []string{"take-screenshot", "s-1"},
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
		test.Fatal("a device capability succeeded without the device helper installed")
	}
	if !strings.Contains(diagnostics.String(), "installed") {
		test.Fatalf("the unprepared capability was not explained to the person: %q", diagnostics.String())
	}
}

func TestDevicesUnprepared(test *testing.T) {
	test.Parallel()

	var diagnostics bytes.Buffer
	process, failure := host.New(host.Options{
		Arguments: []string{"list-devices"},
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
		test.Fatal("a device capability succeeded without the device helper installed")
	}
	if !strings.Contains(diagnostics.String(), "installed") {
		test.Fatalf("the unprepared capability was not explained to the person: %q", diagnostics.String())
	}
}
