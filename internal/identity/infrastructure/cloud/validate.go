package cloud

import "errors"

func (options Options) validate() error {
	switch {
	case options.Agent == nil:
		return errors.New("cloud agent is required")
	case options.Base == "":
		return errors.New("cloud base URL is required")
	case options.Provider != nil && options.Clock == nil:
		return errors.New("cloud clock is required with a provider")
	default:
		return nil
	}
}
