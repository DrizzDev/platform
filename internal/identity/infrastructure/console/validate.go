package console

import "errors"

func (options Options) validate() error {
	if options.Writer == nil {
		return errors.New("console writer is required")
	}
	return nil
}
