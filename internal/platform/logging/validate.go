package logging

import "errors"

func (options Options) validate() error {
	switch {
	case options.Output == nil:
		return errors.New("logging output is required")
	case options.Build.Name() == "":
		return errors.New("logging build identity is required")
	default:
		return nil
	}
}
