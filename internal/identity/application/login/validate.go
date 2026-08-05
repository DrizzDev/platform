package login

import "errors"

func (options Options) validate() error {
	switch {
	case options.Establishment == nil:
		return errors.New("login establishment is required")
	case options.Publication == nil:
		return errors.New("login publication is required")
	case options.Authority == nil:
		return errors.New("login authority is required")
	case options.Clock == nil:
		return errors.New("login clock is required")
	default:
		return nil
	}
}
