package control

import (
	"context"
	"errors"

	"github.com/DrizzDev/platform/internal/device/domain/code"
	"github.com/DrizzDev/platform/internal/device/domain/device"
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

func (flow Flow) Act(scope context.Context, contact touch.Touch) Action {
	if failure := flow.bridge.Tap(scope, contact); failure != nil {
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
