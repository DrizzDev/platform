package platform_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/platform"
)

func TestPlatform(test *testing.T) {
	test.Parallel()

	for _, known := range []platform.Platform{platform.Android, platform.Simulator, platform.Handset} {
		if !known.Valid() {
			test.Fatalf("platform %q rejected", known)
		}
	}
	if platform.Platform("WINDOWS").Valid() {
		test.Fatal("an unknown platform was accepted")
	}
	if platform.Android.String() != "ANDROID" {
		test.Fatalf("platform string = %q", platform.Android.String())
	}
}
