package control

import (
	"context"

	"github.com/DrizzDev/platform/internal/device/domain/capture"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/touch"
)

// Discoverer lists the devices currently reachable through the bridge.
type Discoverer interface {
	List(context.Context) ([]device.Device, error)
}

// Observer captures the current screen of a target device.
type Observer interface {
	Screenshot(context.Context, device.Device) (capture.Capture, error)
}

// Actor performs a single input gesture on a target device.
type Actor interface {
	Tap(context.Context, touch.Touch) error
}

// Bridge is the neutral device-adapter port. The infrastructure adapter over the device-bridge sidecar implements it;
// the flow depends only on this contract, never on a vendor or transport type. New capabilities extend one role
// interface without disturbing the others.
type Bridge interface {
	Actor
	Observer
	Discoverer
}
