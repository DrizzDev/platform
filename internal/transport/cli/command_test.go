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

func (performer performer) Snapshot(context.Context, operator.Target) (operator.Snapshot, error) {
	return operator.Snapshot{}, performer.fail
}

func (performer performer) Hierarchy(context.Context, operator.Target) (operator.Tree, error) {
	return operator.Tree{}, performer.fail
}

func (performer performer) Dimensions(context.Context, operator.Target) (operator.Extent, error) {
	return operator.Extent{}, performer.fail
}

func (performer performer) Devices(context.Context) (operator.Roster, error) {
	return performer.roster, performer.fail
}

func (performer performer) Install(context.Context, operator.Package) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Launch(context.Context, operator.Application) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Terminate(context.Context, operator.Application) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Wipe(context.Context, operator.Application) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Installed(context.Context, operator.Target) (operator.Listing, error) {
	return operator.Listing{}, performer.fail
}

func (performer performer) Running(context.Context, operator.Target) (operator.Listing, error) {
	return operator.Listing{}, performer.fail
}

func (performer performer) Foreground(context.Context, operator.Target) (operator.Report, error) {
	return operator.Report{}, performer.fail
}

func (performer performer) Url(context.Context, operator.Target) (operator.Report, error) {
	return operator.Report{}, performer.fail
}

func (performer performer) Disk(context.Context, operator.Target) (operator.Measure, error) {
	return operator.Measure{}, performer.fail
}

func (performer performer) Images(context.Context) (operator.Images, error) {
	return operator.Images{}, performer.fail
}

func (performer performer) Boot(context.Context, operator.Image) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Pause(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Resume(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Tap(context.Context, operator.Contact) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Swipe(context.Context, operator.Drag) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Pinch(context.Context, operator.Squeeze) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Press(context.Context, operator.Key) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Type(context.Context, operator.Entry) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Clear(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Back(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Home(context.Context, operator.Target) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
}

func (performer performer) Locate(context.Context, operator.Fix) (operator.Ack, error) {
	return operator.Ack{}, performer.fail
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
		arguments: []string{"take-screenshot", "s-1"},
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
		arguments: []string{"list-devices"},
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
		arguments: []string{"take-screenshot", "s-9"},
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
