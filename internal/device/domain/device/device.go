package device

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/platform"
	"github.com/DrizzDev/platform/internal/device/domain/serial"
)

// Device is the immutable identity of one target: the family it belongs to and the serial that addresses it.
type Device struct {
	serial   serial.Serial
	platform platform.Platform
}

type Input struct {
	Serial   serial.Serial
	Platform platform.Platform
}

func New(input Input) (Device, error) {
	device := Device{platform: input.Platform, serial: input.Serial}
	if failure := device.validate(); failure != nil {
		return Device{}, failure
	}
	return device, nil
}

func (device Device) Platform() platform.Platform {
	return device.platform
}

func (device Device) Serial() serial.Serial {
	return device.serial
}

func (device Device) validate() error {
	switch {
	case !device.platform.Valid():
		return errors.New("device platform is invalid")
	case device.serial.String() == "":
		return errors.New("device serial is required")
	default:
		return nil
	}
}
