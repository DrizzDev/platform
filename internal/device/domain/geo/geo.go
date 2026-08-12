package geo

import (
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/device"
)

const (
	latitude  = 90.0
	longitude = 180.0
)

// Fix is a location to report on a device: a latitude and longitude in degrees.
type Fix struct {
	device device.Device
	lat    float64
	lon    float64
}

type Input struct {
	Device device.Device
	Lat    float64
	Lon    float64
}

func New(input Input) (Fix, error) {
	made := Fix{device: input.Device, lat: input.Lat, lon: input.Lon}
	if failure := made.validate(); failure != nil {
		return Fix{}, failure
	}
	return made, nil
}

func (fix Fix) Device() device.Device {
	return fix.device
}

func (fix Fix) Lat() float64 {
	return fix.lat
}

func (fix Fix) Lon() float64 {
	return fix.lon
}

func (fix Fix) validate() error {
	switch {
	case fix.device.Serial().String() == "":
		return errors.New("location device is required")
	case fix.lat < -latitude || fix.lat > latitude:
		return errors.New("latitude is out of range")
	case fix.lon < -longitude || fix.lon > longitude:
		return errors.New("longitude is out of range")
	default:
		return nil
	}
}
