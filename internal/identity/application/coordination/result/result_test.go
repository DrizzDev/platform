package result_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/result"
)

func TestValid(test *testing.T) {
	test.Parallel()

	for _, value := range []result.Result{result.Published, result.Rejected, result.Uncertain} {
		if !value.Valid() {
			test.Fatalf("result %q was rejected", value)
		}
		if value.String() != string(value) {
			test.Fatalf("result = %q", value.String())
		}
	}
}

func TestUnknown(test *testing.T) {
	test.Parallel()

	if result.Result("OTHER").Valid() {
		test.Fatal("an unknown result was accepted")
	}
}
