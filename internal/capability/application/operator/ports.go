package operator

import (
	"context"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/device"
)

// flow performs device work behind the neutral device port; the device control flow satisfies it.
type flow interface {
	Observe(context.Context, device.Device) control.Observation
	Discover(context.Context) control.Discovery
}

// recorder opens one execution record; the capture recorder satisfies it. Recording is observational, so a failure to
// open or write a record never fails the capability it observes.
type recorder interface {
	Begin() (*recording.Execution, error)
}
