package code_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/code"
)

func TestCode(test *testing.T) {
	test.Parallel()

	known := []code.Code{
		code.Missing, code.Timeout, code.Failed, code.Protected,
		code.Cancelled, code.Unavailable, code.Incompatible, code.Unauthorized,
	}
	for _, value := range known {
		if !value.Valid() {
			test.Fatalf("code %q rejected", value)
		}
		if value.Detail() == "" {
			test.Fatalf("code %q has no agent detail", value)
		}
	}
	if code.Code("DEVICE_ON_FIRE").Valid() {
		test.Fatal("an unknown code was accepted")
	}
}

func TestRetryable(test *testing.T) {
	test.Parallel()

	if !code.Timeout.Retryable() || !code.Unavailable.Retryable() {
		test.Fatal("a transient code was reported non-retryable")
	}
	if code.Missing.Retryable() || code.Protected.Retryable() {
		test.Fatal("a terminal code was reported retryable")
	}
}
