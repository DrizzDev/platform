package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/transport/cli"
)

type connector struct {
	captured bool
}

func (*connector) Survey(context.Context) connect.Report { return connect.Report{} }

func (*connector) Enable(context.Context, connect.Selection) (connect.Report, error) {
	return connect.Report{}, nil
}

func (*connector) Disable(context.Context, connect.Selection) connect.Report { return connect.Report{} }

func (fake *connector) Capture(context.Context, connect.Selection) (connect.Report, error) {
	fake.captured = true
	return connect.Report{}, nil
}

func (*connector) Uncapture(context.Context, connect.Selection) connect.Report {
	return connect.Report{}
}

const disclosure = "Recording your prompts and responses for context"

func TestConnectEnableDiscloses(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	fake := &connector{}
	options := fixture{arguments: []string{"connect", "enable", "claude-code"}, output: &output}.options()
	options.Connect = fake

	command, failure := cli.New(options)
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if !strings.Contains(output.String(), disclosure) {
		test.Fatalf("enable did not disclose capture: %q", output.String())
	}
	if !fake.captured {
		test.Fatal("enable did not record context by default")
	}
}

func TestConnectEnableNoCapture(test *testing.T) {
	test.Parallel()

	var output bytes.Buffer
	fake := &connector{}
	options := fixture{arguments: []string{"connect", "enable", "claude-code", "--no-capture"}, output: &output}.options()
	options.Connect = fake

	command, failure := cli.New(options)
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Run(context.Background()); failure != nil {
		test.Fatal(failure)
	}
	if strings.Contains(output.String(), disclosure) {
		test.Fatalf("--no-capture still disclosed capture: %q", output.String())
	}
	if fake.captured {
		test.Fatal("--no-capture still recorded context")
	}
}
