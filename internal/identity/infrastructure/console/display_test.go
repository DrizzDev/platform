package console_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/device"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/console"
)

func TestShow(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	display, failure := console.New(console.Options{Writer: &output})
	if failure != nil {
		test.Fatal(failure)
	}
	instruction := device.Instruction{User: "WDJB-MJHT", Verification: "https://drizz.dev/activate"}
	if failure := display.Show(context.Background(), instruction); failure != nil {
		test.Fatal(failure)
	}
	text := output.String()
	if !strings.Contains(text, "WDJB-MJHT") || !strings.Contains(text, "https://drizz.dev/activate") {
		test.Fatalf("output = %q", text)
	}
}

func TestWriterRequired(test *testing.T) {
	test.Parallel()

	if _, failure := console.New(console.Options{}); failure == nil {
		test.Fatal("a missing writer was accepted")
	}
}
