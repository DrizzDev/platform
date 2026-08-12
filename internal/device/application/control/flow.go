package control

import (
	"context"
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/bundle"
	"github.com/DrizzDev/platform/internal/device/domain/code"
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

// Flow drives a device through the neutral bridge port and maps every bridge outcome to a stable, agent-facing code.
// It never inspects a vendor or sidecar detail. New capabilities are added as further methods over the same port.
type Flow struct {
	bridge Bridge
}

func (flow Flow) Discover(scope context.Context) Discovery {
	devices, failure := flow.bridge.List(scope)
	if failure != nil {
		return Discovery{outcome: flow.deny(failure)}
	}
	return Discovery{devices: devices}
}

func (flow Flow) Observe(scope context.Context, target device.Device) Observation {
	frame, failure := flow.bridge.Screenshot(scope, target)
	if failure != nil {
		return Observation{outcome: flow.deny(failure)}
	}
	return Observation{capture: frame}
}

func (flow Flow) Snapshot(scope context.Context, target device.Device) Portrait {
	frame, tree, failure := flow.bridge.Snapshot(scope, target)
	if failure != nil {
		return Portrait{outcome: flow.deny(failure)}
	}
	return Portrait{capture: frame, hierarchy: tree}
}

func (flow Flow) Hierarchy(scope context.Context, target device.Device) Reading {
	tree, failure := flow.bridge.Hierarchy(scope, target)
	if failure != nil {
		return Reading{outcome: flow.deny(failure)}
	}
	return Reading{text: tree}
}

func (flow Flow) Dimensions(scope context.Context, target device.Device) Extent {
	width, height, failure := flow.bridge.Dimensions(scope, target)
	if failure != nil {
		return Extent{outcome: flow.deny(failure)}
	}
	return Extent{width: width, height: height}
}

func (flow Flow) Act(scope context.Context, contact touch.Touch) Action {
	if failure := flow.bridge.Tap(scope, contact); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Swipe(scope context.Context, drag swipe.Swipe) Action {
	if failure := flow.bridge.Swipe(scope, drag); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Pinch(scope context.Context, gesture pinch.Pinch) Action {
	if failure := flow.bridge.Pinch(scope, gesture); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Press(scope context.Context, key press.Press) Action {
	if failure := flow.bridge.Press(scope, key); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Type(scope context.Context, entry text.Text) Action {
	if failure := flow.bridge.Type(scope, entry); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Clear(scope context.Context, target device.Device) Action {
	if failure := flow.bridge.Clear(scope, target); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Back(scope context.Context, target device.Device) Action {
	if failure := flow.bridge.Back(scope, target); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Home(scope context.Context, target device.Device) Action {
	if failure := flow.bridge.Home(scope, target); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Locate(scope context.Context, fix geo.Fix) Action {
	if failure := flow.bridge.Locate(scope, fix); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Install(scope context.Context, payload parcel.Parcel) Action {
	if failure := flow.bridge.Install(scope, payload); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Launch(scope context.Context, application bundle.Bundle) Action {
	if failure := flow.bridge.Launch(scope, application); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Terminate(scope context.Context, application bundle.Bundle) Action {
	if failure := flow.bridge.Terminate(scope, application); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Wipe(scope context.Context, application bundle.Bundle) Action {
	if failure := flow.bridge.Wipe(scope, application); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Installed(scope context.Context, target device.Device) Inventory {
	apps, failure := flow.bridge.Installed(scope, target)
	if failure != nil {
		return Inventory{outcome: flow.deny(failure)}
	}
	return Inventory{apps: apps}
}

func (flow Flow) Running(scope context.Context, target device.Device) Inventory {
	apps, failure := flow.bridge.Running(scope, target)
	if failure != nil {
		return Inventory{outcome: flow.deny(failure)}
	}
	return Inventory{apps: apps}
}

func (flow Flow) Foreground(scope context.Context, target device.Device) Reading {
	name, failure := flow.bridge.Foreground(scope, target)
	if failure != nil {
		return Reading{outcome: flow.deny(failure)}
	}
	return Reading{text: name}
}

func (flow Flow) Url(scope context.Context, target device.Device) Reading {
	link, failure := flow.bridge.Url(scope, target)
	if failure != nil {
		return Reading{outcome: flow.deny(failure)}
	}
	return Reading{text: link}
}

func (flow Flow) Disk(scope context.Context, target device.Device) Measure {
	value, failure := flow.bridge.Disk(scope, target)
	if failure != nil {
		return Measure{outcome: flow.deny(failure)}
	}
	return Measure{value: value}
}

func (flow Flow) Images(scope context.Context, family platform.Platform) Catalogue {
	names, failure := flow.bridge.Images(scope, family)
	if failure != nil {
		return Catalogue{outcome: flow.deny(failure)}
	}
	return Catalogue{names: names}
}

func (flow Flow) Boot(scope context.Context, image emulator.Boot) Action {
	if failure := flow.bridge.Boot(scope, image); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Pause(scope context.Context, target device.Device) Action {
	if failure := flow.bridge.Pause(scope, target); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) Resume(scope context.Context, target device.Device) Action {
	if failure := flow.bridge.Resume(scope, target); failure != nil {
		return Action{outcome: flow.deny(failure)}
	}
	return Action{}
}

func (flow Flow) deny(cause error) outcome {
	kind := code.Unavailable
	var carrier interface{ reason() code.Code }

	switch {
	case errors.Is(cause, context.Canceled), errors.Is(cause, context.DeadlineExceeded):
		kind = code.Cancelled
	case errors.As(cause, &carrier):
		kind = carrier.reason()
	}
	return outcome{failure: kind}
}
