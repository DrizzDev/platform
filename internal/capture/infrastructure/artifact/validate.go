package artifact

import "errors"

func (options Options) validate() error {
	switch {
	case options.Root == "":
		return errors.New("artifact root is required")
	case options.Ceiling < 0:
		return errors.New("artifact ceiling must not be negative")
	case options.Logger == nil:
		return errors.New("artifact logger is required")
	case options.Tracer == nil:
		return errors.New("artifact tracer is required")
	case options.Meter == nil:
		return errors.New("artifact meter is required")
	default:
		return nil
	}
}
