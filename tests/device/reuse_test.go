//go:build device

package device_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/device/infrastructure/bridge"
)

// TestReuse checks that one driver keeps the helper alive across calls: a second call should be far faster than the
// first (which pays the process start-up), proving the helper is reused rather than respawned per request.
func TestReuse(test *testing.T) {
	location := os.Getenv("DRIZZ_DEVICE_SIDECAR")
	digest := os.Getenv("DRIZZ_DEVICE_DIGEST")
	if location == "" || digest == "" {
		test.Skip("set DRIZZ_DEVICE_SIDECAR and DRIZZ_DEVICE_DIGEST to run the reuse check")
	}

	driver, failure := bridge.New(bridge.Options{Location: location, Digest: digest, Timeout: 30 * time.Second})
	if failure != nil {
		test.Fatal(failure)
	}
	defer func() { _ = driver.Close() }()

	scope := context.Background()

	start := time.Now()
	if _, failure := driver.List(scope); failure != nil {
		test.Fatal(failure)
	}
	first := time.Since(start)

	start = time.Now()
	if _, failure := driver.List(scope); failure != nil {
		test.Fatal(failure)
	}
	second := time.Since(start)

	test.Logf("first=%v second=%v", first, second)
	if second > first/2 {
		test.Fatalf("second call was not substantially faster (first=%v second=%v): helper was respawned, not reused", first, second)
	}
}
