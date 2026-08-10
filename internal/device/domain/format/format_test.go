package format_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/format"
)

func TestFormat(test *testing.T) {
	test.Parallel()

	if !format.Png.Valid() {
		test.Fatal("png rejected")
	}
	if format.Format("GIF").Valid() {
		test.Fatal("an unsupported format was accepted")
	}
	if format.Png.String() != "PNG" {
		test.Fatalf("format string = %q", format.Png.String())
	}
}
