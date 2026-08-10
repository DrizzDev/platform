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
	default:
		return nil
	}
}
