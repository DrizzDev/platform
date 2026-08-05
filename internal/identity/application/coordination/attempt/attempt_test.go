package attempt_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/attempt"
	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
)

func TestAttempt(test *testing.T) {
	test.Parallel()

	value, failure := attempt.New(attempt.Input{Revision: 1, Epoch: epoch.Epoch(3)})
	if failure != nil {
		test.Fatal(failure)
	}
	if value.Revision() != 1 || value.Epoch() != epoch.Epoch(3) {
		test.Fatalf("attempt = %+v", value)
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	if _, failure := attempt.New(attempt.Input{Revision: 0}); failure == nil {
		test.Fatal("an attempt without a revision was accepted")
	}
}
