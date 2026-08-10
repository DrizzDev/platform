package mark_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/mark"
)

func TestMark(test *testing.T) {
	test.Parallel()

	if !mark.Exact.Valid() || !mark.Inferred.Valid() {
		test.Fatal("a known mark was rejected")
	}
	if mark.Mark("MAYBE").Valid() {
		test.Fatal("an unknown mark was accepted")
	}
	if mark.Inferred.String() != "INFERRED" {
		test.Fatalf("mark string = %q", mark.Inferred.String())
	}
}
