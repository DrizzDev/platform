package vault

import "errors"

func (options Options) validate() error {
	if options.Store == nil {
		return errors.New("vault store is required")
	}
	return nil
}
