package point_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/point"
)

func TestPoint(test *testing.T) {
	test.Parallel()

	value, failure := point.New(point.Input{X: 540, Y: 1200})
	if failure != nil {
		test.Fatal(failure)
	}
	if value.X() != 540 || value.Y() != 1200 {
		test.Fatalf("point = %d,%d", value.X(), value.Y())
	}
}

func TestPointRejects(test *testing.T) {
	test.Parallel()

	rejected := []point.Input{{X: -1, Y: 0}, {X: 0, Y: -1}, {X: 1 << 17, Y: 0}}
	for _, sample := range rejected {
		if _, failure := point.New(sample); failure == nil {
			test.Fatalf("point %v was accepted", sample)
		}
	}
}
