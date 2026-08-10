package recording

import "errors"

func (options Options) validate() error {
	switch {
	case options.Writer == nil:
		return errors.New("record writer is required")
	case options.Sink == nil:
		return errors.New("record sink is required")
	case options.Logger == nil:
		return errors.New("record logger is required")
	case options.Keeper == nil:
		return errors.New("record keeper is required")
	case options.Clock == nil:
		return errors.New("record clock is required")
	case options.Lease <= 0:
		return errors.New("record lease must be positive")
	default:
		return nil
	}
}
