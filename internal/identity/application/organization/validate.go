package organization

import "errors"

func (options Options) validate() error {
	if options.Resolver == nil {
		return errors.New("organization resolver is required")
	}
	return nil
}
