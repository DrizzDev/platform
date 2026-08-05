package device

import "errors"

func (options Options) validate() error {
	switch {
	case options.Provider == nil:
		return errors.New("device provider is required")
	case options.Display == nil:
		return errors.New("device display is required")
	default:
		return nil
	}
}
