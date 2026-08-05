package method_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/method"
)

func TestValid(test *testing.T) {
	test.Parallel()

	for _, value := range []method.Method{method.Browser, method.Device, method.Workload} {
		if !value.Valid() {
			test.Fatalf("method %q was rejected", value)
		}
		if value.String() != string(value) {
			test.Fatalf("method = %q", value.String())
		}
	}
}

func TestUnknown(test *testing.T) {
	test.Parallel()

	if method.Method("OTHER").Valid() {
		test.Fatal("an unknown method was accepted")
	}
}
