package binary_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/integration/infrastructure/binary"
)

func TestLocateReturnsPath(test *testing.T) {
	test.Parallel()

	path, failure := binary.New().Locate()
	if failure != nil {
		test.Fatal(failure)
	}
	if path == "" {
		test.Fatal("the resolved executable path must not be empty")
	}
}
