package failure

import "errors"

const limit = 256

func (value Value) validate() error {
	switch {
	case !value.code.Valid():
		return errors.New("failure code is invalid")
	case len(value.detail) > limit:
		return errors.New("failure detail exceeds the maximum length")
	case len(value.correlation) > limit:
		return errors.New("failure correlation exceeds the maximum length")
	case value.retry < 0:
		return errors.New("failure retry must not be negative")
	default:
		return nil
	}
}
