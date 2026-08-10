package courier

import "errors"

func (options Options) validate() error {
	switch {
	case options.Ledger == nil:
		return errors.New("courier ledger is required")
	case options.Vault == nil:
		return errors.New("courier vault is required")
	case options.Uploader == nil:
		return errors.New("courier uploader is required")
	default:
		return nil
	}
}
