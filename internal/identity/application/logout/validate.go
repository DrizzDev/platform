package logout

import "errors"

func (options Options) validate() error {
	switch {
	case options.Vault == nil:
		return errors.New("logout vault is required")
	case options.Publication == nil:
		return errors.New("logout publication is required")
	case options.Revocation == nil:
		return errors.New("logout revocation is required")
	case options.Clock == nil:
		return errors.New("logout clock is required")
	default:
		return nil
	}
}
