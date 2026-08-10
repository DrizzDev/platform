package bearings_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
)

func TestBearings(test *testing.T) {
	test.Parallel()

	session, failure := identifier.New("session-1")
	if failure != nil {
		test.Fatal(failure)
	}
	value := bearings.New(bearings.Input{Session: session})
	if value.Session().String() != "session-1" {
		test.Fatalf("session = %q", value.Session().String())
	}
	if !value.Turn().Empty() || !value.Capability().Empty() {
		test.Fatal("an absent dimension is not reported empty")
	}
}
