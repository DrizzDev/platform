package trace_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

func TestTrace(test *testing.T) {
	test.Parallel()

	value, failure := trace.New("01HEXECUTION")
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != "01HEXECUTION" {
		test.Fatalf("trace = %q", value.String())
	}
	if _, failure := trace.New(""); failure == nil {
		test.Fatal("an empty trace was accepted")
	}
}
