package control

import "github.com/DrizzDev/platform/internal/device/domain/device"

// Discovery is the reachable device set, or a code-only failure.
type Discovery struct {
	outcome
	devices []device.Device
}

func (discovery Discovery) Devices() []device.Device {
	return append([]device.Device(nil), discovery.devices...)
}
