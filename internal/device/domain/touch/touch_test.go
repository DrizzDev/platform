package touch_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
	"github.com/DrizzDev/platform/internal/device/domain/point"
	"github.com/DrizzDev/platform/internal/device/domain/serial"
	"github.com/DrizzDev/platform/internal/device/domain/touch"
)

func (fixture fixture) target(test *testing.T) device.Device {
	test.Helper()
	handle, failure := serial.New("emulator-5554")
	if failure != nil {
		test.Fatal(failure)
	}
	target, failure := device.New(device.Input{Platform: platform.Android, Serial: handle})
	if failure != nil {
		test.Fatal(failure)
	}
	return target
}

type fixture struct{}

func TestTouch(test *testing.T) {
	test.Parallel()

	spot, failure := point.New(point.Input{X: 540, Y: 1200})
	if failure != nil {
		test.Fatal(failure)
	}
	contact, failure := touch.New(touch.Input{Device: fixture{}.target(test), Point: spot})
	if failure != nil {
		test.Fatal(failure)
	}
	if contact.Point().X() != 540 || contact.Device().Serial().String() != "emulator-5554" {
		test.Fatal("touch did not carry its device and point")
	}
}

func TestTouchRejectsMissingDevice(test *testing.T) {
	test.Parallel()

	if _, failure := touch.New(touch.Input{}); failure == nil {
		test.Fatal("a touch without a device was accepted")
	}
}
