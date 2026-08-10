package touch

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/point"
)

// Touch is a single tap: the target device and the coordinate to press.
type Touch struct {
	point  point.Point
	device device.Device
}

type Input struct {
	Point  point.Point
	Device device.Device
}

func New(input Input) (Touch, error) {
	touch := Touch{device: input.Device, point: input.Point}
	if failure := touch.validate(); failure != nil {
		return Touch{}, failure
	}
	return touch, nil
}

func (touch Touch) Device() device.Device {
	return touch.device
}

func (touch Touch) Point() point.Point {
	return touch.point
}

func (touch Touch) validate() error {
	if touch.device.Serial().String() == "" {
		return errors.New("touch device is required")
	}
	return nil
}
