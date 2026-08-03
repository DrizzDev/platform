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

	default:
		return nil
	}
}
