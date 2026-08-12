package bundle

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/device"
)

// Bundle names one application on a device: the device and the application's package or bundle identifier.
type Bundle struct {
	device device.Device
	id     string
}

type Input struct {
	Device device.Device
	Id     string
}

func New(input Input) (Bundle, error) {
	made := Bundle{device: input.Device, id: input.Id}
	if failure := made.validate(); failure != nil {
		return Bundle{}, failure
	}
	return made, nil
}

func (bundle Bundle) Device() device.Device {
	return bundle.device
}

func (bundle Bundle) Id() string {
	return bundle.id
}

func (bundle Bundle) validate() error {
	switch {
	case bundle.device.Serial().String() == "":
		return errors.New("bundle device is required")
	case bundle.id == "":
		return errors.New("bundle identifier is required")
	default:
		return nil
	}
}
