package operator

import (
	"context"
	"fmt"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
	"github.com/DrizzDev/platform/internal/device/domain/capture"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/serial"
)

// Screenshot captures the screen of the target device and records the capture. It resolves the serial to a connected
// device first, so a missing or unauthorized device is reported as a stable code rather than a raw failure.
func (operator Operator) Screenshot(scope context.Context, target Target) (shot Shot, failure error) {
	scope, watch := operator.begin(scope, catalog.Screenshot)
	defer func() { watch.finish(scope, failure) }()

	subject, failure := operator.resolve(scope, target.Serial)
	if failure != nil {
		return Shot{}, failure
	}
	observation := operator.flow.Observe(scope, subject)
	if reason, failed := observation.Failure(); failed {
		return Shot{}, operator.refuse(reason)
	}
	frame := observation.Capture()
	operator.record(scope, frame)
	return Shot{Image: frame.Image().Bytes(), Format: frame.Format().String()}, nil
}

func (operator Operator) resolve(scope context.Context, value string) (device.Device, error) {
	chosen, failure := serial.New(value)
	if failure != nil {
		return device.Device{}, Refusal{Code: outcome.Invalid}
	}
	discovery := operator.flow.Discover(scope)
	if reason, failed := discovery.Failure(); failed {
		return device.Device{}, operator.refuse(reason)
	}
	for _, candidate := range discovery.Devices() {
		if candidate.Serial().String() == chosen.String() {
			return candidate, nil
		}
	}
	return device.Device{}, Refusal{Code: outcome.Missing}
}

// record writes the capture as one observational execution entry, with the screen image as its artifact.
func (operator Operator) record(scope context.Context, frame capture.Capture) {
	operator.inscribe(scope, entry{
		capability: catalog.Screenshot,
		note: recording.Note{
			Origin:   origin.Capability,
			Fidelity: fidelity.Exact,
			Category: category.Screen,
			Payload:  []byte(fmt.Sprintf("%dx%d", frame.Width(), frame.Height())),
			Artifact: frame.Image().Bytes(),
		},
	})
}
