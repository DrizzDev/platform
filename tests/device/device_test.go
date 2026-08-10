//go:build device

package device_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/device/application/control"
	"github.com/DrizzDev/platform/internal/device/domain/point"
	"github.com/DrizzDev/platform/internal/device/domain/touch"
	"github.com/DrizzDev/platform/internal/device/infrastructure/bridge"
)

// TestJourney drives a real Android target through the Go core and the reused device-bridge sidecar — discover a
// device, capture its screen, tap the centre — running only under the `device` tag when a sidecar binary, its digest,
// and a connected device are supplied, so it stays out of the general gate.
func TestJourney(test *testing.T) {

	digest := os.Getenv("DRIZZ_DEVICE_DIGEST")
	location := os.Getenv("DRIZZ_DEVICE_SIDECAR")

	if location == "" || digest == "" {
		test.Skip("set DRIZZ_DEVICE_SIDECAR and DRIZZ_DEVICE_DIGEST to run the device journey")
	}

	driver, failure := bridge.New(bridge.Options{Location: location, Digest: digest, Timeout: 30 * time.Second})
	if failure != nil {
		test.Fatal(failure)
	}
	defer func() { _ = driver.Close() }()

	flow, failure := control.New(control.Options{Bridge: driver})
	if failure != nil {
		test.Fatal(failure)
	}
	scope := context.Background()

	discovery := flow.Discover(scope)
	if discovery.Failed() {
		reason, _ := discovery.Failure()
		test.Fatalf("discover: %s", reason)
	}
	if len(discovery.Devices()) == 0 {
		test.Skip("no device connected")
	}
	target := discovery.Devices()[0]

	observation := flow.Observe(scope, target)
	if observation.Failed() {
		reason, _ := observation.Failure()
		test.Fatalf("observe: %s", reason)
	}
	shot := observation.Capture()
	if shot.Width() <= 0 || shot.Height() <= 0 || shot.Image().Size() == 0 {
		test.Fatalf("empty capture: %dx%d/%d", shot.Width(), shot.Height(), shot.Image().Size())
	}

	spot, failure := point.New(point.Input{X: shot.Width() / 2, Y: shot.Height() / 2})
	if failure != nil {
		test.Fatal(failure)
	}
	contact, failure := touch.New(touch.Input{Device: target, Point: spot})
	if failure != nil {
		test.Fatal(failure)
	}
	if action := flow.Act(scope, contact); action.Failed() {
		reason, _ := action.Failure()
		test.Fatalf("tap: %s", reason)
	}
	test.Logf("device journey ok: %s %dx%d", target.Serial(), shot.Width(), shot.Height())
}
