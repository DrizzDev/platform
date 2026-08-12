package mcp

import "errors"

func (options Options) validate() error {
	switch {
	case options.Logger == nil:
		return errors.New("MCP logger is required")
	case options.External == nil:
		return errors.New("MCP external logger is required")
	case options.Tracer == nil:
		return errors.New("MCP tracer is required")
	case options.Meter == nil:
		return errors.New("MCP meter is required")
	case options.Input == nil:
		return errors.New("MCP input is required")
	case options.Output == nil:
		return errors.New("MCP output is required")
	case options.Perform == nil:
		return errors.New("MCP capability operator is required")
	case options.Release.Name() == "":
		return errors.New("MCP release identity is required")
	default:
		return nil
	}
}
