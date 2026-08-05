package browser

import (
	"errors"
	"strings"
)

func (options Options) validate() error {
	switch {
	case options.Opener == nil:
		return errors.New("browser opener is required")
	case !strings.HasPrefix(options.Address, "127.0.0.1:"):
		return errors.New("browser address must bind loopback only")
	case options.Path == "":
		return errors.New("browser callback path is required")
	default:
		return nil
	}
}
