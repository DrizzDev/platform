package cli_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/application/release"
	"github.com/DrizzDev/platform/internal/transport/cli"
)

type runner func(context.Context) error

type fixture struct {
	arguments []string
	output    io.Writer
	server    cli.Runner
}

func (run runner) Run(scope context.Context) error {
	return run(scope)
}

func TestVersion(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	command, failure := cli.New(fixture{
		arguments: []string{"version"},
		output:    &output,
		server: runner(func(context.Context) error {
			test.Fatal("MCP called")
			return nil
		}),
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if !strings.HasPrefix(output.String(), "drizz ") {
		test.Fatalf("output = %q", output.String())
	}
}

func TestMCP(test *testing.T) {
	test.Parallel()

	called := false
	command, failure := cli.New(fixture{
		arguments: []string{"mcp"},
		output:    io.Discard,
		server: runner(func(context.Context) error {
			called = true
			return nil
		}),
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if !called {
		test.Fatal("MCP was not called")
	}
}

func TestDependencies(test *testing.T) {
	test.Parallel()

	_, failure := cli.New(cli.Options{})
	if failure == nil {
		test.Fatal("missing dependencies were accepted")
	}
}

func (fixture fixture) options() cli.Options {
	identity, _ := release.New(release.Input{
		Name:     "drizz",
		Version:  "1.0.0",
		Revision: "revision_123",
	})
	return cli.Options{
		Arguments: fixture.arguments,
		Streams: cli.Streams{
			Input:   strings.NewReader(""),
			Output:  fixture.output,
			Failure: io.Discard,
		},
		Release: identity,
		MCP:     fixture.server,
	}
}
