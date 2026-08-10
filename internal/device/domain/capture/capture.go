package capture

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/format"
	"github.com/DrizzDev/platform/internal/device/domain/image"
)

const limit = 1 << 16

// Capture is one observation of a device screen: the encoded frame and the
// dimensions an agent reasons over.
type Capture struct {
	width  int
	height int
	image  image.Image
	format format.Format
}

type Input struct {
	Width  int
	Height int
	Image  image.Image
	Format format.Format
}

func New(input Input) (Capture, error) {
	capture := Capture{format: input.Format, image: input.Image, width: input.Width, height: input.Height}
	if failure := capture.validate(); failure != nil {
		return Capture{}, failure
	}
	return capture, nil
}

func (capture Capture) Format() format.Format {
	return capture.format
}

func (capture Capture) Image() image.Image {
	return capture.image
}

func (capture Capture) Width() int {
	return capture.width
}

func (capture Capture) Height() int {
	return capture.height
}

func (capture Capture) validate() error {
	switch {
	case !capture.format.Valid():
		return errors.New("capture format is invalid")
	case capture.image.Size() == 0:
		return errors.New("capture image is required")
	case capture.width <= 0 || capture.height <= 0:
		return errors.New("capture dimensions must be positive")
	case capture.width > limit || capture.height > limit:
		return errors.New("capture dimensions exceed the maximum")
	default:
		return nil
	}
}
