package cli_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/application/release"
	"github.com/DrizzDev/platform/internal/capability/application/operator"
	"github.com/DrizzDev/platform/internal/capability/domain/outcome"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/application/logout"
	"github.com/DrizzDev/platform/internal/transport/cli"
)

type runner func(context.Context) error

type authenticator func(context.Context) (login.Result, error)

type exit func(context.Context) (logout.Result, error)

type fixture struct {
	arguments []string
	output    io.Writer
	server    cli.Runner
	session   cli.Session
	device    cli.Session
	logout    cli.Departure
	perform   cli.Perform
}

type performer struct {
	shot   operator.Shot
	roster operator.Roster
	fail   error
}

func (performer performer) Screenshot(context.Context, operator.Target) (operator.Shot, error) {
	return performer.shot, performer.fail
}

func (performer performer) Devices(context.Context) (operator.Roster, error) {
	return performer.roster, performer.fail
}

func (run runner) Run(scope context.Context) error {
	return run(scope)
}

func (run authenticator) Run(scope context.Context) (login.Result, error) {
	return run(scope)
}

func (run exit) Run(scope context.Context) (logout.Result, error) {
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

func TestLogin(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	command, failure := cli.New(fixture{
		arguments: []string{"login"},
		output:    &output,
		server:    runner(func(context.Context) error { test.Fatal("MCP called"); return nil }),
		session:   authenticator(func(context.Context) (login.Result, error) { return login.Result{}, nil }),
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if output.String() != "Signed in to Drizz.\n" {
		test.Fatalf("output = %q", output.String())
	}
}

func TestLoginUnavailable(test *testing.T) {
	test.Parallel()

	command, failure := cli.New(fixture{
		arguments: []string{"login"},
		output:    io.Discard,
		server:    runner(func(context.Context) error { test.Fatal("MCP called"); return nil }),
		session: authenticator(func(context.Context) (login.Result, error) {
			return login.Result{}, errors.New("composition failed")
		}),
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	outcome := command.Run(context.Background())
	if outcome == nil || outcome.Error() != "Drizz authentication is temporarily unavailable. Try again." {
		test.Fatalf("outcome = %v", outcome)
	}
}

func TestUsage(test *testing.T) {
	test.Parallel()

	command, failure := cli.New(fixture{
		arguments: []string{"login", "unexpected"},
		output:    io.Discard,
		server:    runner(func(context.Context) error { test.Fatal("MCP called"); return nil }),
		session: authenticator(func(context.Context) (login.Result, error) {
			test.Fatal("login ran on invalid syntax")
			return login.Result{}, nil
		}),
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	outcome := command.Run(context.Background())
	if outcome == nil || outcome.Error() != "Usage: drizz login [--device]" {
		test.Fatalf("outcome = %v", outcome)
	}
}

func TestDevice(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	command, failure := cli.New(fixture{
		arguments: []string{"login", "--device"},
		output:    &output,
		server:    runner(func(context.Context) error { test.Fatal("MCP called"); return nil }),
		session: authenticator(func(context.Context) (login.Result, error) {
			test.Fatal("browser login ran for --device")
			return login.Result{}, nil
		}),
		device: authenticator(func(context.Context) (login.Result, error) { return login.Result{}, nil }),
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if output.String() != "Signed in to Drizz.\n" {
		test.Fatalf("output = %q", output.String())
	}
}

func TestSignout(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	command, failure := cli.New(fixture{
		arguments: []string{"logout"},
		output:    &output,
		server:    runner(func(context.Context) error { test.Fatal("MCP called"); return nil }),
		logout:    exit(func(context.Context) (logout.Result, error) { return logout.Result{}, nil }),
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if output.String() != "Signed out of Drizz.\n" {
		test.Fatalf("output = %q", output.String())
	}
}

func TestScreenshot(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	command, failure := cli.New(fixture{
		arguments: []string{"screenshot", "s-1"},
		output:    &output,
		perform:   performer{shot: operator.Shot{Image: []byte("png-bytes"), Format: "PNG"}},
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	path := strings.TrimSpace(output.String())
	if !strings.HasSuffix(path, ".png") {
		test.Fatalf("printed path = %q", path)
	}
	written, failure := os.ReadFile(path)
	if failure != nil {
		test.Fatal(failure)
	}
	_ = os.Remove(path)
	if string(written) != "png-bytes" {
		test.Fatalf("written = %q", written)
	}
}

func TestDevices(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	command, failure := cli.New(fixture{
		arguments: []string{"devices"},
		output:    &output,
		perform:   performer{roster: operator.Roster{Serials: []string{"s-1", "s-2"}}},
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if !strings.Contains(output.String(), "s-1") || !strings.Contains(output.String(), "s-2") {
		test.Fatalf("output = %q", output.String())
	}
}

func TestScreenshotRefused(test *testing.T) {
	test.Parallel()

	command, failure := cli.New(fixture{
		arguments: []string{"screenshot", "s-9"},
		output:    io.Discard,
		perform:   performer{fail: operator.Refusal{Code: outcome.Missing}},
	}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	failure = command.Run(context.Background())
	if failure == nil || !strings.Contains(failure.Error(), "not found") {
		test.Fatalf("refused screenshot = %v", failure)
	}
}

func (fixture fixture) options() cli.Options {
	identity, _ := release.New(release.Input{
		Name:     "drizz",
		Version:  "1.0.0",
		Revision: "revision_123",
	})
	session := fixture.session
	if session == nil {
		session = authenticator(func(context.Context) (login.Result, error) { return login.Result{}, nil })
	}
	passcode := fixture.device
	if passcode == nil {
		passcode = authenticator(func(context.Context) (login.Result, error) { return login.Result{}, nil })
	}
	farewell := fixture.logout
	if farewell == nil {
		farewell = exit(func(context.Context) (logout.Result, error) { return logout.Result{}, nil })
	}
	perform := fixture.perform
	if perform == nil {
		perform = performer{}
	}
	server := fixture.server
	if server == nil {
		server = runner(func(context.Context) error { return nil })
	}
	return cli.Options{
		Arguments: fixture.arguments,
		Streams: cli.Streams{
			Input:   strings.NewReader(""),
			Output:  fixture.output,
			Failure: io.Discard,
		},
		Release: identity,
		MCP:     server,
		Login:   session,
		Device:  passcode,
		Logout:  farewell,
		Perform: perform,
	}
}
