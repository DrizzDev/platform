package janitor

import "errors"

func (options Options) validate() error {
	switch {
	case options.Archive == nil:
		return errors.New("janitor archive is required")
	case options.Vault == nil:
		return errors.New("janitor vault is required")
	case options.Notifier == nil:
		return errors.New("janitor notifier is required")
	case options.Clock == nil:
		return errors.New("janitor clock is required")
	case options.Ceiling < 0 || options.Relief < 0 || options.Grace < 0:
		return errors.New("janitor budget and grace must not be negative")
	case options.low() >= options.budget():
		return errors.New("janitor relief must stay below the ceiling")
	default:
		return nil
	}
}
