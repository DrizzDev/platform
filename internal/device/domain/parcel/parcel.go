package parcel

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/device"
)

// Parcel is an application package to install on a device: the device and the path to the package file.
type Parcel struct {
	device device.Device
	path   string
}

type Input struct {
	Device device.Device
	Path   string
}

func New(input Input) (Parcel, error) {
	made := Parcel{device: input.Device, path: input.Path}
	if failure := made.validate(); failure != nil {
		return Parcel{}, failure
	}
	return made, nil
}

func (parcel Parcel) Device() device.Device {
	return parcel.device
}

func (parcel Parcel) Path() string {
	return parcel.path
}

func (parcel Parcel) validate() error {
	switch {
	case parcel.device.Serial().String() == "":
		return errors.New("parcel device is required")
	case parcel.path == "":
		return errors.New("parcel path is required")
	default:
		return nil
	}
}
