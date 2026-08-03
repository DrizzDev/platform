package build_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/platform/build"
)

func TestInfo(test *testing.T) {
	test.Parallel()

	info := build.Read()
	if info.Name() != "drizz" {
		test.Fatalf("name = %q", info.Name())
	}
	if info.Version() == "" {
		test.Fatal("version is empty")
	}
	if info.Revision() == "" {
		test.Fatal("revision is empty")
	}
}
