// Package server models the stdio MCP server entry the installer writes into an agent application's configuration:
// the name it is filed under and the command plus arguments that launch it. The value is agent-neutral; the installer
// builds one Drizz entry and every dialect writer renders it in its own format.
package server

import "errors"

// Server is one stdio MCP server entry: the name it is filed under in an agent's configuration, and the command and
// arguments that start it over standard input and output.
type Server struct {
	name    string
	command string
	args    []string
}

// Input constructs a Server. Name and Command are required; a server with no name to file it under or no command to
// launch is not a usable entry.
type Input struct {
	Name    string
	Command string
	Args    []string
}

func New(input Input) (Server, error) {
	if input.Name == "" {
		return Server{}, errors.New("server name is required")
	}
	if input.Command == "" {
		return Server{}, errors.New("server command is required")
	}
	return Server{
		name:    input.Name,
		command: input.Command,
		args:    append([]string(nil), input.Args...),
	}, nil
}

func (server Server) Name() string {
	return server.name
}

func (server Server) Command() string {
	return server.command
}

func (server Server) Args() []string {
	return append([]string(nil), server.args...)
}
