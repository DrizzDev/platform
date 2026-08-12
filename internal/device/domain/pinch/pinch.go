package pinch

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/point"
)

// Pinch is a two-finger zoom around a centre point: the radius grows to zoom in and shrinks to zoom out.
type Pinch struct {
	device device.Device
	centre point.Point
	from   int
	to     int
}

type Input struct {
	Device device.Device
	Centre point.Point
	From   int
	To     int
}

func New(input Input) (Pinch, error) {
	made := Pinch{device: input.Device, centre: input.Centre, from: input.From, to: input.To}
	if failure := made.validate(); failure != nil {
		return Pinch{}, failure
	}
	return made, nil
}

func (pinch Pinch) Device() device.Device {
	return pinch.device
}

func (pinch Pinch) Centre() point.Point {
	return pinch.centre
}

func (pinch Pinch) From() int {
	return pinch.from
}

func (pinch Pinch) To() int {
	return pinch.to
}

func (pinch Pinch) validate() error {
	switch {
	case pinch.device.Serial().String() == "":
		return errors.New("pinch device is required")
	case pinch.from < 0 || pinch.to < 0:
		return errors.New("pinch radius must not be negative")
	default:
		return nil
	}
}
