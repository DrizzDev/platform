package control_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/capture"
	"github.com/DrizzDev/platform/internal/device/domain/code"
	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/format"
	"github.com/DrizzDev/platform/internal/device/domain/image"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
	"github.com/DrizzDev/platform/internal/device/domain/serial"
	"github.com/DrizzDev/platform/internal/device/domain/touch"
)

// bridge is a scripted Bridge: it returns its configured result, or its
// configured failure when one is set.
type bridge struct {
	devices []device.Device
	capture capture.Capture
	failure error
}

func (bridge bridge) List(context.Context) ([]device.Device, error) {
	return bridge.devices, bridge.failure
}

func (bridge bridge) Screenshot(context.Context, device.Device) (capture.Capture, error) {
	return bridge.capture, bridge.failure
}

func (bridge bridge) Tap(context.Context, touch.Touch) error {
	return bridge.failure
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) flow(stub bridge) control.Flow {
	fixture.test.Helper()
	flow, failure := control.New(control.Options{Bridge: stub})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return flow
}

func (fixture fixture) target() device.Device {
	fixture.test.Helper()
	handle, failure := serial.New("emulator-5554")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	target, failure := device.New(device.Input{Platform: platform.Android, Serial: handle})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return target
}

func (fixture fixture) frame() capture.Capture {
	fixture.test.Helper()
	pixels, failure := image.New([]byte{1, 2, 3})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	shot, failure := capture.New(capture.Input{Format: format.Png, Image: pixels, Width: 1080, Height: 2400})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return shot
}

func TestDiscover(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	flow := kit.flow(bridge{devices: []device.Device{kit.target()}})
	result := flow.Discover(context.Background())
	if result.Failed() {
		test.Fatalf("discover failed: %v", result)
	}
	if len(result.Devices()) != 1 || result.Devices()[0].Serial().String() != "emulator-5554" {
		test.Fatalf("discover returned %d devices", len(result.Devices()))
	}
}

func TestObserve(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	flow := kit.flow(bridge{capture: kit.frame()})
	result := flow.Observe(context.Background(), kit.target())
	if result.Failed() {
		test.Fatalf("observe failed: %v", result)
	}
	if result.Capture().Width() != 1080 {
		test.Fatalf("observe width = %d", result.Capture().Width())
	}
}

func TestAct(test *testing.T) {
	test.Parallel()

	flow := fixture{test: test}.flow(bridge{})
	result := flow.Act(context.Background(), touch.Touch{})
	if result.Failed() {
		test.Fatalf("act failed: %v", result)
	}
}

func TestFailures(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	cases := []struct {
		cause error
		want  code.Code
	}{
		{control.Missing{}, code.Missing},
		{control.Unauthorized{}, code.Unauthorized},
		{control.Protected{}, code.Protected},
		{control.Timeout{}, code.Timeout},
		{control.Incompatible{}, code.Incompatible},
		{control.Failed{}, code.Failed},
		{control.Unavailable{}, code.Unavailable},
		{context.Canceled, code.Cancelled},
		{context.DeadlineExceeded, code.Cancelled},
		{errors.New("unmapped"), code.Unavailable},
	}
	for _, sample := range cases {
		flow := kit.flow(bridge{failure: sample.cause})
		result := flow.Act(context.Background(), touch.Touch{})
		got, present := result.Failure()
		if !present || got != sample.want {
			test.Fatalf("cause %v mapped to %q, want %q", sample.cause, got, sample.want)
		}
	}
}
