package connect

import "errors"

func (options Options) validate() error {
	switch {
	case options.Resolver == nil:
		return errors.New("installer resolver is required")
	case options.Store == nil:
		return errors.New("installer store is required")
	case options.Recorder == nil:
		return errors.New("installer recorder is required")
	case options.Monitor == nil:
		return errors.New("installer monitor is required")
	case options.Logger == nil:
		return errors.New("installer logger is required")
	default:
		return nil
	}
}
