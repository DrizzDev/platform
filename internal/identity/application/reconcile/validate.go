package reconcile

import "errors"

func (options Options) validate() error {
	switch {
	case options.Queue == nil:
		return errors.New("reconcile queue is required")
	case options.Vault == nil:
		return errors.New("reconcile vault is required")
	case options.Clock == nil:
		return errors.New("reconcile clock is required")
	default:
		return nil
	}
}
