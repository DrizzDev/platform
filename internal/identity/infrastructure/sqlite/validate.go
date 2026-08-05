package sqlite

import "errors"

func (options Options) validate() error {
	switch {
	case options.Path == "":
		return errors.New("database path is required")
	case options.Logger == nil:
		return errors.New("database logger is required")
	case options.Tracer == nil:
		return errors.New("database tracer is required")
	case options.Meter == nil:
		return errors.New("database meter is required")
	default:
		return nil
	}
}
