package swipe

import (
	"errors"
	"time"

	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/point"
)

// Swipe is a drag across the screen: the device, the start and end points, and how long the gesture takes.
type Swipe struct {
	device device.Device
	from   point.Point
	to     point.Point
	span   time.Duration
}

type Input struct {
	Device device.Device
	From   point.Point
	To     point.Point
	Span   time.Duration
}

func New(input Input) (Swipe, error) {
	made := Swipe{device: input.Device, from: input.From, to: input.To, span: input.Span}
	if failure := made.validate(); failure != nil {
		return Swipe{}, failure
	}
	return made, nil
}

func (swipe Swipe) Device() device.Device {
	return swipe.device
}

func (swipe Swipe) From() point.Point {
	return swipe.from
}

func (swipe Swipe) To() point.Point {
	return swipe.to
}

func (swipe Swipe) Span() time.Duration {
	return swipe.span
}

func (swipe Swipe) validate() error {
	if swipe.device.Serial().String() == "" {
		return errors.New("swipe device is required")
	}
	if swipe.span < 0 {
		return errors.New("swipe duration must not be negative")
	}
	return nil
}
