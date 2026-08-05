package cli

import "errors"

func (options Options) validate() error {
	switch {
	case options.Streams.Input == nil:
		return errors.New("CLI input is required")

	case options.Streams.Output == nil:
		return errors.New("CLI output is required")

	case options.Streams.Failure == nil:
		return errors.New("CLI failure output is required")

	case options.Release.Name() == "":
		return errors.New("CLI release identity is required")

	case options.MCP == nil:
		return errors.New("CLI MCP server is required")

	case options.Login == nil:
		return errors.New("CLI login session is required")

	case options.Device == nil:
		return errors.New("CLI device session is required")

	case options.Logout == nil:
		return errors.New("CLI logout is required")

	default:
		return nil
	}
}
