package device

import "github.com/DrizzDev/platform/internal/identity/application/login"

func New(options Options) (login.Establishment, error) {
	if failure := options.validate(); failure != nil {
		return nil, failure
	}
	return terminal{provider: options.Provider, display: options.Display}, nil
}
