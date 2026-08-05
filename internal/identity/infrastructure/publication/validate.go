package publication

import "errors"

func (options Options) validate() error {
	switch {
	case options.Vault == nil:
		return errors.New("publication vault is required")
	case options.Ledger == nil:
		return errors.New("publication ledger is required")
	case options.Random == nil:
		return errors.New("publication random is required")
	case options.Session.String() == "":
		return errors.New("publication session is required")
	default:
		return nil
	}
}
