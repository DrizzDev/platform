package image

import "errors"

const limit = 32 << 20

// Image is the raw bytes of one captured frame.
// It defensively clones on the way in and out so the value stays immutable.
type Image struct {
	data []byte
}

func New(data []byte) (Image, error) {
	image := Image{data: append([]byte(nil), data...)}
	if failure := image.validate(); failure != nil {
		return Image{}, failure
	}
	return image, nil
}

func (image Image) Bytes() []byte {
	return append([]byte(nil), image.data...)
}

func (image Image) Size() int {
	return len(image.data)
}

func (image Image) validate() error {
	switch {
	case len(image.data) == 0:
		return errors.New("image data is required")
	case len(image.data) > limit:
		return errors.New("image data exceeds the maximum size")
	default:
		return nil
	}
}
