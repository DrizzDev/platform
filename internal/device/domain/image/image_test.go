package image_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/image"
)

func TestImage(test *testing.T) {
	test.Parallel()

	source := []byte{1, 2, 3, 4}
	value, failure := image.New(source)
	if failure != nil {
		test.Fatal(failure)
	}
	if value.Size() != 4 {
		test.Fatalf("size = %d", value.Size())
	}

	source[0] = 9
	if value.Bytes()[0] != 1 {
		test.Fatal("image did not clone its source")
	}
	value.Bytes()[1] = 9
	if value.Bytes()[1] != 2 {
		test.Fatal("image did not clone on read")
	}
}

func TestImageRejectsEmpty(test *testing.T) {
	test.Parallel()

	if _, failure := image.New(nil); failure == nil {
		test.Fatal("empty image was accepted")
	}
}
