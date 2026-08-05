package failure_test

import (
	"strings"
	"testing"
	"time"

	fault "github.com/DrizzDev/platform/internal/identity/application/failure"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/action"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/category"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

func TestValue(test *testing.T) {
	test.Parallel()

	value, failure := fault.New(fault.Input{
		Code:        code.Unavailable,
		Detail:      "the identity provider did not respond",
		Correlation: "correlation_123",
		Retry:       2 * time.Second,
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if value.Code() != code.Unavailable {
		test.Fatalf("code = %q", value.Code())
	}
	if value.Category() != category.Authentication || value.Action() != action.Retry || !value.Retryable() {
		test.Fatalf("derived attributes = %q %q %v", value.Category(), value.Action(), value.Retryable())
	}
	if value.Retry() != 2*time.Second || value.Correlation() != "correlation_123" {
		test.Fatalf("value = %+v", value)
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	cases := map[string]fault.Input{
		"code":        {Code: code.Code("OTHER")},
		"detail":      {Code: code.Failed, Detail: strings.Repeat("a", 257)},
		"correlation": {Code: code.Failed, Correlation: strings.Repeat("a", 257)},
		"retry":       {Code: code.Failed, Retry: -1},
	}
	for name, input := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			if _, failure := fault.New(input); failure == nil {
				test.Fatal("invalid failure value was accepted")
			}
		})
	}
}
