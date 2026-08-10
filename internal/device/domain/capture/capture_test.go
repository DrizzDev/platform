package capture_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/capture"
	"github.com/DrizzDev/platform/internal/device/domain/format"
	"github.com/DrizzDev/platform/internal/device/domain/image"
)

func (fixture fixture) frame(test *testing.T) image.Image {
	test.Helper()
	frame, failure := image.New([]byte{1, 2, 3})
	if failure != nil {
		test.Fatal(failure)
	}
	return frame
}

type fixture struct{}

func TestCapture(test *testing.T) {
	test.Parallel()

	shot, failure := capture.New(capture.Input{
		Format: format.Png,
		Image:  fixture{}.frame(test),
		Width:  1080,
		Height: 2400,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if shot.Width() != 1080 || shot.Height() != 2400 || shot.Format() != format.Png {
		test.Fatalf("capture = %dx%d %q", shot.Width(), shot.Height(), shot.Format())
	}
}

func TestCaptureRejects(test *testing.T) {
	test.Parallel()

	frame := fixture{}.frame(test)
	rejected := map[string]capture.Input{
		"format":     {Format: format.Format("GIF"), Image: frame, Width: 10, Height: 10},
		"dimensions": {Format: format.Png, Image: frame, Width: 0, Height: 10},
	}
	for name, sample := range rejected {
		if _, failure := capture.New(sample); failure == nil {
			test.Fatalf("%s capture was accepted", name)
		}
	}
}
