package cloud

import "errors"

func (options Options) validate() error {
	switch {
	case options.Agent == nil:
		return errors.New("cloud http client is required")
	case options.Provider == nil:
		return errors.New("cloud credential provider is required")
	case options.Clock == nil:
		return errors.New("cloud clock is required")
	case options.Base == "":
		return errors.New("cloud base url is required")
	default:
		return nil
	}
}
