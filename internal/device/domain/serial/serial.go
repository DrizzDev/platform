package serial

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const limit = 256

// Serial is a device's transport identity — the value adb or the sidecar uses to address one device.
type Serial struct {
	value string
}

func New(value string) (Serial, error) {
	serial := Serial{value: value}
	if failure := serial.validate(); failure != nil {
		return Serial{}, failure
	}
	return serial, nil
}

func (serial Serial) String() string {
	return serial.value
}

func (serial Serial) validate() error {
	switch {
	case serial.value == "":
		return errors.New("device serial is required")
	case len(serial.value) > limit:
		return errors.New("device serial exceeds the maximum length")
	case !utf8.ValidString(serial.value):
		return errors.New("device serial is not valid UTF-8")
	case strings.ContainsFunc(serial.value, unicode.IsControl):
		return errors.New("device serial contains a control character")
	default:
		return nil
	}
}
