package serial_test

import (
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/device/domain/serial"
)

func TestSerial(test *testing.T) {
	test.Parallel()

	value, failure := serial.New("emulator-5554")
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != "emulator-5554" {
		test.Fatalf("serial = %q", value.String())
	}
}

func TestSerialRejects(test *testing.T) {
	test.Parallel()

	rejected := map[string]string{
		"empty":   "",
		"control": "abc\x00def",
		"long":    strings.Repeat("a", 257),
	}
	for name, value := range rejected {
		if _, failure := serial.New(value); failure == nil {
			test.Fatalf("%s serial was accepted", name)
		}
	}
}
