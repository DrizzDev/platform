package session

import "errors"

func (options Options) validate() error {
	switch {
	case options.Vault == nil:
		return errors.New("session vault is required")
	case options.Refresh == nil:
		return errors.New("session refresh is required")
	case options.Publication == nil:
		return errors.New("session publication is required")
	case options.Epoch == nil:
		return errors.New("session epoch is required")
	case options.Clock == nil:
		return errors.New("session clock is required")
	default:
		return nil
	}
}
