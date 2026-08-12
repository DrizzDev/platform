package operator

import (
	"context"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/bundle"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/emulator"
	"github.com/DrizzDev/platform/internal/device/domain/geo"
	"github.com/DrizzDev/platform/internal/device/domain/parcel"
	"github.com/DrizzDev/platform/internal/device/domain/pinch"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
	"github.com/DrizzDev/platform/internal/device/domain/press"
	"github.com/DrizzDev/platform/internal/device/domain/swipe"
	"github.com/DrizzDev/platform/internal/device/domain/text"
	"github.com/DrizzDev/platform/internal/device/domain/touch"
)

// flow performs device work behind the neutral device port; the device control flow satisfies it. Each capability
// uses one method; new capabilities extend this port.
type flow interface {
	Observe(context.Context, device.Device) control.Observation
	Snapshot(context.Context, device.Device) control.Portrait
	Hierarchy(context.Context, device.Device) control.Reading
	Dimensions(context.Context, device.Device) control.Extent
	Discover(context.Context) control.Discovery
	Act(context.Context, touch.Touch) control.Action
	Install(context.Context, parcel.Parcel) control.Action
	Launch(context.Context, bundle.Bundle) control.Action
	Terminate(context.Context, bundle.Bundle) control.Action
	Wipe(context.Context, bundle.Bundle) control.Action
	Installed(context.Context, device.Device) control.Inventory
	Running(context.Context, device.Device) control.Inventory
	Foreground(context.Context, device.Device) control.Reading
	Url(context.Context, device.Device) control.Reading
	Disk(context.Context, device.Device) control.Measure
	Images(context.Context, platform.Platform) control.Catalogue
	Boot(context.Context, emulator.Boot) control.Action
	Pause(context.Context, device.Device) control.Action
	Resume(context.Context, device.Device) control.Action
	Swipe(context.Context, swipe.Swipe) control.Action
	Pinch(context.Context, pinch.Pinch) control.Action
	Press(context.Context, press.Press) control.Action
	Type(context.Context, text.Text) control.Action
	Clear(context.Context, device.Device) control.Action
	Back(context.Context, device.Device) control.Action
	Home(context.Context, device.Device) control.Action
	Locate(context.Context, geo.Fix) control.Action
}

// recorder opens one execution record; the capture recorder satisfies it. Recording is observational, so a failure to
// open or write a record never fails the capability it observes.
type recorder interface {
	Begin() (*recording.Execution, error)
}
