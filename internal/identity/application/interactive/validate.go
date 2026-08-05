package interactive

import "errors"

func (options Options) validate() error {
	switch {
	case options.Authorization == nil:
		return errors.New("interactive authorization is required")
	case options.Browser == nil:
		return errors.New("interactive browser is required")
	case options.Random == nil:
		return errors.New("interactive random is required")
	default:
		return nil
	}
}
