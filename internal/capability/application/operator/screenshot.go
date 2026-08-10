package operator

import (
	"context"
	"fmt"
	"log/slog"

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
func (operator Operator) Screenshot(scope context.Context, target Target) (Shot, error) {
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
		return device.Device{}, Refusal{code: outcome.Invalid}
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
	return device.Device{}, Refusal{code: outcome.Missing}
}

// record writes the capture as one observational execution entry. If the record cannot be opened, the drop is logged
// and swallowed, so recording can never break the capability it observes but a lost record is never silent.
func (operator Operator) record(scope context.Context, frame capture.Capture) {
	execution, failure := operator.recorder.Begin()
	if failure != nil {
		operator.logger.WarnContext(scope, "capability.record.dropped", slog.String("capability", catalog.Screenshot))
		return
	}
	execution.Record(scope, recording.Note{
		Origin:   origin.Capability,
		Fidelity: fidelity.Exact,
		Category: category.Screen,
		Payload:  []byte(fmt.Sprintf("%dx%d", frame.Width(), frame.Height())),
		Artifact: frame.Image().Bytes(),
	})
}
