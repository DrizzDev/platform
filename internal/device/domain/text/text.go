package text

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/device"
)

// ceiling bounds one typed string so a single request stays small and predictable.
const ceiling = 4096

// Text is a string to type into the focused field on a device.
type Text struct {
	device  device.Device
	content string
}

type Input struct {
	Device  device.Device
	Content string
}

func New(input Input) (Text, error) {
	made := Text{device: input.Device, content: input.Content}
	if failure := made.validate(); failure != nil {
		return Text{}, failure
	}
	return made, nil
}

func (text Text) Device() device.Device {
	return text.device
}

func (text Text) Content() string {
	return text.content
}

func (text Text) validate() error {
	switch {
	case text.device.Serial().String() == "":
		return errors.New("text device is required")
	case text.content == "":
		return errors.New("text content is required")
	case len(text.content) > ceiling:
		return errors.New("text content exceeds the limit")
	default:
		return nil
	}
}
