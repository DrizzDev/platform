package epoch_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/epoch"
)

func TestNext(test *testing.T) {
	test.Parallel()

	if epoch.Epoch(1).Next() != 2 {
		test.Fatal("epoch did not advance")
	}
}

func TestOrder(test *testing.T) {
	test.Parallel()

	if !epoch.Epoch(1).Before(2) || !epoch.Epoch(2).After(1) {
		test.Fatal("epoch ordering is wrong")
	}
	if epoch.Epoch(2).Before(2) {
		test.Fatal("an equal epoch must not compare before")
	}
}
