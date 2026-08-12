package intake

import "errors"

func (options Options) validate() error {
	switch {
	case options.Recorder == nil:
		return errors.New("intake recorder is required")
	case options.Monitor == nil:
		return errors.New("intake monitor is required")
	case options.Logger == nil:
		return errors.New("intake logger is required")
	default:
		return nil
	}
}
