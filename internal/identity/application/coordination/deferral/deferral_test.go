package deferral_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/deferral"
)

func TestDeferral(test *testing.T) {
	test.Parallel()

	postponed, failure := deferral.New(deferral.Input{Key: "key-0", Next: time.Unix(1000, 0)})
	if failure != nil {
		test.Fatal(failure)
	}
	if postponed.Key() != "key-0" {
		test.Fatalf("deferral = %+v", postponed)
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	if _, failure := deferral.New(deferral.Input{Next: time.Unix(1000, 0)}); failure == nil {
		test.Fatal("a deferral without a key was accepted")
	}
}
