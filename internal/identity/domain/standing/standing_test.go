package standing_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/standing"
)

func TestValid(test *testing.T) {
	test.Parallel()

	for _, value := range []standing.Standing{standing.Active, standing.Expired, standing.Revoked} {
		if !value.Valid() {
			test.Fatalf("standing %q was rejected", value)
		}
		if value.String() != string(value) {
			test.Fatalf("standing = %q", value.String())
		}
	}
}

func TestUnknown(test *testing.T) {
	test.Parallel()

	if standing.Standing("OTHER").Valid() {
		test.Fatal("an unknown standing was accepted")
	}
}
