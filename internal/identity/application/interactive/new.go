package interactive

import "github.com/DrizzDev/platform/internal/identity/application/login"

func New(options Options) (login.Establishment, error) {
	if failure := options.validate(); failure != nil {
		return nil, failure
	}
	return front{
		authorization: options.Authorization,
		browser:       options.Browser,
		random:        options.Random,
	}, nil
}
