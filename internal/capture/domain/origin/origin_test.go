package origin_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/origin"
)

func TestOrigin(test *testing.T) {
	test.Parallel()

	if !origin.Capability.Valid() || !origin.Host.Valid() {
		test.Fatal("a known origin was rejected")
	}
	if origin.Origin("CLOUD").Valid() {
		test.Fatal("an unknown origin was accepted")
	}
	if origin.Capability.String() != "CAPABILITY" {
		test.Fatalf("origin string = %q", origin.Capability.String())
	}
}
