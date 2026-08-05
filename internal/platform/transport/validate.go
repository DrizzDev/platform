package transport

import "errors"

func (options Options) validate() error {
	switch {
	case options.Timeout <= 0:
		return errors.New("transport timeout is required")
	case options.Dial <= 0:
		return errors.New("transport dial timeout is required")
	case options.Ceiling <= 0:
		return errors.New("transport body ceiling is required")
	case options.Retries < 0:
		return errors.New("transport retries must not be negative")
	case options.Tracing == nil:
		return errors.New("transport tracing provider is required")
	case options.Metering == nil:
		return errors.New("transport metering provider is required")
	default:
		return nil
	}
}
