package operator

import "errors"

func (options Options) validate() error {
	switch {
	case options.Flow == nil:
		return errors.New("operator flow is required")
	case options.Recorder == nil:
		return errors.New("operator recorder is required")
	case options.Logger == nil:
		return errors.New("operator logger is required")
	default:
		return nil
	}
}
