package host

import (
	"github.com/DrizzDev/platform/internal/application/release"
	"github.com/DrizzDev/platform/internal/platform/build"
	"github.com/DrizzDev/platform/internal/transport/cli"
)

func New(options Options) (*Host, error) {
	if failure := options.validate(); failure != nil {
		return nil, failure
	}
	info := build.Read()
	identity, failure := release.New(release.Input{
		Name:     info.Name(),
		Version:  info.Version(),
		Revision: info.Revision(),
	})
	if failure != nil {
		return nil, failure
	}
	base := foundation{environment: options.Environment, streams: options.Streams, build: info}
	device := &station{base: base}
	command, failure := cli.New(cli.Options{
		Arguments: options.Arguments,
		Streams: cli.Streams{
			Input:   options.Streams.Input,
			Output:  options.Streams.Output,
			Failure: options.Streams.Failure,
		},
		Release: identity,
		MCP: runtime{
			environment: options.Environment,
			streams:     options.Streams,
			identity:    identity,
			build:       info,
		},
		Login:    access{base},
		Device:   terminal{base},
		Logout:   departure{base},
		Perform:  pilot{station: device},
		Connect:  connector{base: base},
		Receiver: receiver{base: base},
	})
	if failure != nil {
		return nil, failure
	}
	return &Host{command: command, failure: options.Streams.Failure, device: device}, nil
}
