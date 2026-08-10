package span_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/span"
)

func TestSpan(test *testing.T) {
	test.Parallel()

	value, failure := span.New("01HHOP")
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != "01HHOP" {
		test.Fatalf("span = %q", value.String())
	}
	if !(span.Span{}).Empty() {
		test.Fatal("the zero span is not reported empty (needed to mark the root)")
	}
	if _, failure := span.New(""); failure == nil {
		test.Fatal("an empty span was accepted")
	}
}
