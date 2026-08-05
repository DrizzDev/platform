package pointer_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/pointer"
)

func TestPointer(test *testing.T) {
	test.Parallel()

	value, failure := pointer.New(pointer.Input{Key: "session_123#1", Revision: 1, Epoch: epoch.Epoch(2)})
	if failure != nil {
		test.Fatal(failure)
	}
	if value.Key() != "session_123#1" || value.Revision() != 1 || value.Epoch() != epoch.Epoch(2) {
		test.Fatalf("pointer = %+v", value)
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	cases := map[string]pointer.Input{
		"key":      {Revision: 1},
		"revision": {Key: "session_123#1"},
	}
	for name, input := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			if _, failure := pointer.New(input); failure == nil {
				test.Fatal("invalid pointer was accepted")
			}
		})
	}
}
