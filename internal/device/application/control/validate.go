package control

import "errors"

func (options Options) validate() error {
	if options.Bridge == nil {
		return errors.New("device bridge is required")
	}
	return nil
}
