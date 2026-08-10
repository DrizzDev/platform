package recording

import "errors"

func (options Options) validate() error {
	switch {
	case options.Writer == nil:
		return errors.New("record writer is required")
	case options.Sink == nil:
		return errors.New("record sink is required")
	case options.Logger == nil:
		return errors.New("record logger is required")
	default:
		return nil
	}
}
