package bridge

import "errors"

func (options Options) validate() error {
	switch {
	case options.Location == "":
		return errors.New("device sidecar location is required")
	case options.Digest == "":
		return errors.New("device sidecar digest is required")
	case options.Timeout <= 0:
		return errors.New("device request timeout is required")
	case options.Tracer == nil:
		return errors.New("device sidecar tracer is required")
	case options.Meter == nil:
		return errors.New("device sidecar meter is required")
	default:
		return nil
	}
}
