package lobby

import "errors"

func (options Options) validate() error {
	switch {
	case options.Register == nil:
		return errors.New("lobby register is required")
	case options.Clock == nil:
		return errors.New("lobby clock is required")
	case options.Window < 0:
		return errors.New("lobby window must not be negative")
	case options.Capacity < 0:
		return errors.New("lobby capacity must not be negative")
	default:
		return nil
	}
}
