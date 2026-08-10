package fidelity_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
)

func TestFidelity(test *testing.T) {
	test.Parallel()

	known := []fidelity.Fidelity{fidelity.Exact, fidelity.Summary, fidelity.Inferred, fidelity.Unavailable, fidelity.Redacted}
	for _, value := range known {
		if !value.Valid() {
			test.Fatalf("fidelity %q rejected", value)
		}
	}
	if fidelity.Fidelity("GUESS").Valid() {
		test.Fatal("an unknown fidelity was accepted")
	}
	if fidelity.Exact.String() != "EXACT" {
		test.Fatalf("fidelity string = %q", fidelity.Exact.String())
	}
}
