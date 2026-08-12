package mcp

import (
	protocol "github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel/metric"
)

func New(options Options) (Server, error) {
	if failure := options.validate(); failure != nil {
		return Server{}, failure
	}
	duration, failure := options.Meter.Float64Histogram(
		"drizz.operation.duration",
		metric.WithUnit("s"),
	)
	if failure != nil {
		return Server{}, failure
	}
	implementation := &protocol.Implementation{
		Title:       "Drizz",
		Name:        options.Release.Name(),
		Version:     options.Release.Version(),
		Description: "Local Drizz capabilities for agents and developer tools.",
	}
	server := protocol.NewServer(implementation, &protocol.ServerOptions{
		Logger:       options.External,
		Capabilities: &protocol.ServerCapabilities{},
		Instructions: "Use these Drizz tools for every device action — seeing the screen, tapping, typing, and managing " +
			"apps and emulators — rather than the device's native command-line tools. Screen and hierarchy captures return " +
			"the path of the saved file alongside the inline result, so there is no need to re-capture with other tooling.",
	})
	drizz := Server{
		server:   server,
		duration: duration,
		logger:   options.Logger,
		tracer:   options.Tracer,
		input:    options.Input,
		output:   options.Output,
	}
	drizz.register(options.Perform)
	return drizz, nil
}
