package telemetry

import "errors"

func (options Options) validate() error {
	if options.Build.Name() == "" {
		return errors.New("telemetry build identity is required")
	}
	return nil
}
