package host

import "errors"

func (options Options) validate() error {
	switch {
	case options.Streams.Input == nil:
		return errors.New("host input is required")
	case options.Streams.Output == nil:
		return errors.New("host output is required")
	case options.Streams.Failure == nil:
		return errors.New("host failure output is required")
	default:
		return nil
	}
}
