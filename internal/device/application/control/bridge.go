package control

import (
	"context"

	"github.com/DrizzDev/platform/internal/device/domain/app"
	"github.com/DrizzDev/platform/internal/device/domain/bundle"
	"github.com/DrizzDev/platform/internal/device/domain/capture"
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

// Discoverer lists the devices currently reachable through the bridge.
type Discoverer interface {
	List(context.Context) ([]device.Device, error)
}

// Observer reads the current state of a target device — its screen and on-screen element tree.
type Observer interface {
	Screenshot(context.Context, device.Device) (capture.Capture, error)
	Snapshot(context.Context, device.Device) (capture.Capture, string, error)
	Hierarchy(context.Context, device.Device) (string, error)
	Dimensions(context.Context, device.Device) (int, int, error)
}

// Actor performs a single input gesture on a target device. Each interaction is one method over the same port.
type Actor interface {
	Tap(context.Context, touch.Touch) error
	Swipe(context.Context, swipe.Swipe) error
	Pinch(context.Context, pinch.Pinch) error
	Press(context.Context, press.Press) error
	Type(context.Context, text.Text) error
	Clear(context.Context, device.Device) error
	Back(context.Context, device.Device) error
	Home(context.Context, device.Device) error
	Locate(context.Context, geo.Fix) error
	Install(context.Context, parcel.Parcel) error
	Launch(context.Context, bundle.Bundle) error
	Terminate(context.Context, bundle.Bundle) error
	Wipe(context.Context, bundle.Bundle) error
}

// Reader reads application and device metadata from a target device.
type Reader interface {
	Installed(context.Context, device.Device) ([]app.App, error)
	Running(context.Context, device.Device) ([]app.App, error)
	Foreground(context.Context, device.Device) (string, error)
	Url(context.Context, device.Device) (string, error)
	Disk(context.Context, device.Device) (int, error)
}

// Provisioner manages emulators: the images available to run, starting one, and pausing or resuming a running one.
type Provisioner interface {
	Images(context.Context, platform.Platform) ([]string, error)
	Boot(context.Context, emulator.Boot) error
	Pause(context.Context, device.Device) error
	Resume(context.Context, device.Device) error
}

// Bridge is the neutral device-adapter port. The infrastructure adapter over the device-bridge sidecar implements it;
// the flow depends only on this contract, never on a vendor or transport type. New capabilities extend one role
// interface without disturbing the others.
type Bridge interface {
	Actor
	Observer
	Discoverer
	Reader
	Provisioner
}
