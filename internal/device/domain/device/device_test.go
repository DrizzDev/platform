package device_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/device"
	"github.com/DrizzDev/platform/internal/device/domain/platform"
	"github.com/DrizzDev/platform/internal/device/domain/serial"
)

func TestDevice(test *testing.T) {
	test.Parallel()

	handle, failure := serial.New("emulator-5554")
	if failure != nil {
		test.Fatal(failure)
	}
	target, failure := device.New(device.Input{Platform: platform.Android, Serial: handle})
	if failure != nil {
		test.Fatal(failure)
	}
	if target.Platform() != platform.Android || target.Serial().String() != "emulator-5554" {
		test.Fatalf("device = %v/%q", target.Platform(), target.Serial().String())
	}
}

func TestDeviceRejects(test *testing.T) {
	test.Parallel()

	handle, failure := serial.New("emulator-5554")
	if failure != nil {
		test.Fatal(failure)
	}
	if _, failure := device.New(device.Input{Platform: platform.Platform("NOPE"), Serial: handle}); failure == nil {
		test.Fatal("an invalid platform was accepted")
	}
	if _, failure := device.New(device.Input{Platform: platform.Android}); failure == nil {
		test.Fatal("a missing serial was accepted")
	}
}
