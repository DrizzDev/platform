package process_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	protocol "github.com/modelcontextprotocol/go-sdk/mcp"
)

const deadline = 30 * time.Second

const readiness = "mcp.started"

var executable string

func TestMain(suite *testing.M) {
	directory, failure := os.MkdirTemp("", "drizz-process")
	if failure != nil {
		fmt.Fprintln(os.Stderr, failure)
		os.Exit(1)
	}
	executable = filepath.Join(directory, "drizz")
	build := exec.CommandContext(context.Background(), "go", "build", "-trimpath", "-o", executable, "./command/drizz")
	build.Dir = filepath.Join("..", "..")
	if output, failure := build.CombinedOutput(); failure != nil {
		fmt.Fprintf(os.Stderr, "build: %v\n%s", failure, output)
		os.Exit(1)
	}
	// These tests observe the server's own diagnostics (its readiness marker and per-request outcomes), so the suite
	// runs the server with logging on; it is off by default.
	if failure := os.Setenv("DRIZZ_LOG_LEVEL", "info"); failure != nil {
		fmt.Fprintln(os.Stderr, failure)
		os.Exit(1)
	}
	code := suite.Run()
	_ = os.RemoveAll(directory)
	os.Exit(code)
}

func TestNegotiation(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	var diagnostics bytes.Buffer
	command := exec.CommandContext(scope, executable, "mcp")
	command.Stderr = &diagnostics

	client := protocol.NewClient(&protocol.Implementation{Name: "test", Version: "1.0.0"}, nil)
	session, failure := client.Connect(scope, &protocol.CommandTransport{Command: command}, nil)
	if failure != nil {
		test.Fatalf("connect: %v", failure)
	}
	if failure := session.Ping(scope, nil); failure != nil {
		test.Fatalf("ping: %v", failure)
	}
	if name := session.InitializeResult().ServerInfo.Name; name != "drizz" {
		test.Fatalf("server = %q", name)
	}
	if failure := session.Close(); failure != nil {
		test.Fatalf("normal termination reported failure: %v", failure)
	}

	report := diagnostics.String()
	if !strings.Contains(report, readiness) {
		test.Fatalf("diagnostics missing from standard error: %q", report)
	}
	if strings.Contains(report, string(filepath.Separator)+"Users"+string(filepath.Separator)) {
		test.Fatalf("diagnostics leaked an absolute developer path: %q", report)
	}
}

func TestReuse(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	printed, failure := exec.CommandContext(scope, executable, "version").Output()
	if failure != nil {
		test.Fatalf("version: %v", failure)
	}

	command := exec.CommandContext(scope, executable, "mcp")
	client := protocol.NewClient(&protocol.Implementation{Name: "test", Version: "1.0.0"}, nil)
	session, failure := client.Connect(scope, &protocol.CommandTransport{Command: command}, nil)
	if failure != nil {
		test.Fatalf("connect: %v", failure)
	}
	served := session.InitializeResult().ServerInfo
	if failure := session.Close(); failure != nil {
		test.Error(failure)
	}

	if !strings.Contains(string(printed), served.Name) || !strings.Contains(string(printed), served.Version) {
		test.Fatalf("version %q does not match the MCP identity %q %q", printed, served.Name, served.Version)
	}
}

func TestShutdown(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	observer := &monitor{marker: readiness, ready: make(chan struct{})}
	var output bytes.Buffer
	command := exec.CommandContext(scope, executable, "mcp")
	command.Stdout = &output
	command.Stderr = observer
	if _, failure := command.StdinPipe(); failure != nil {
		test.Fatal(failure)
	}
	if failure := command.Start(); failure != nil {
		test.Fatal(failure)
	}

	select {
	case <-observer.ready:
	case <-scope.Done():
		test.Fatal("server did not start before the deadline")
	}
	if failure := command.Process.Signal(syscall.SIGTERM); failure != nil {
		test.Fatal(failure)
	}
	checker := program{test: test}
	if code := checker.code(command.Wait()); code != 0 {
		test.Fatalf("signal shutdown exited with %d", code)
	}
	if output.Len() != 0 {
		test.Fatalf("standard output carried non-protocol content: %q", output.String())
	}
	if report := observer.report(); strings.Contains(report, `"level":"ERROR"`) {
		test.Fatalf("signal shutdown produced an error record: %q", report)
	}
}

func TestStartup(test *testing.T) {
	test.Parallel()

	scope, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	var output, diagnostics bytes.Buffer
	command := exec.CommandContext(scope, executable, "mcp")
	command.Env = append(os.Environ(), "DRIZZ_TELEMETRY_EXPORTER=otlp")
	command.Stdout = &output
	command.Stderr = &diagnostics

	checker := program{test: test}
	code := checker.code(command.Run())
	if code == 0 {
		test.Fatal("genuine startup failure exited successfully")
	}
	if output.Len() != 0 {
		test.Fatalf("standard output carried a diagnostic: %q", output.String())
	}
	report := strings.TrimSpace(diagnostics.String())
	if strings.Count(report, "\n") != 0 || report == "" {
		test.Fatalf("expected exactly one diagnostic line, got: %q", report)
	}
}

func TestRejection(test *testing.T) {
	test.Parallel()

	cases := []struct {
		name  string
		input string
	}{
		{name: "garbage", input: "{ this is not a valid message }\n"},
		{name: "version", input: `{"jsonrpc":"9.9","id":1,"method":"ping"}` + "\n"},
		{name: "trailing", input: `{"jsonrpc":"2.0","id":1,"method":"ping"} trailing` + "\n"},
		{name: "truncated", input: `{"jsonrpc":"2.0","id":1,"me`},
		{name: "multiline", input: "{\n\"jsonrpc\":\"2.0\",\n\"id\":1\n}\n"},
		{name: "oversize", input: strings.Repeat("a", (1<<20)+16) + "\n"},
	}
	for _, item := range cases {
		test.Run(item.name, func(test *testing.T) {
			test.Parallel()
			scope, cancel := context.WithTimeout(context.Background(), deadline)
			defer cancel()

			var output, diagnostics bytes.Buffer
			command := exec.CommandContext(scope, executable, "mcp")
			command.Stdin = strings.NewReader(item.input)
			command.Stdout = &output
			command.Stderr = &diagnostics

			checker := program{test: test}
			if code := checker.code(command.Run()); code == 0 {
				test.Fatal("a client-controlled boundary error exited successfully")
			}
			if output.Len() != 0 {
				test.Fatalf("standard output carried non-protocol content: %q", output.String())
			}
			report := diagnostics.String()
			if strings.Contains(report, "command.failed") || strings.Contains(report, `"level":"ERROR"`) {
				test.Fatalf("client-controlled input was reported as an internal failure: %q", report)
			}
			if !strings.Contains(report, `"outcome":"REJECTED"`) {
				test.Fatalf("client-controlled input was not classified as rejected: %q", report)
			}
		})
	}
}

type program struct {
	test *testing.T
}

func (program program) code(failure error) int {
	if failure == nil {
		return 0
	}
	var exit *exec.ExitError
	if errors.As(failure, &exit) {
		return exit.ExitCode()
	}
	program.test.Fatalf("waiting on the process failed: %v", failure)
	return -1
}

type monitor struct {
	marker string
	buffer strings.Builder
	guard  sync.Mutex
	ready  chan struct{}
	once   sync.Once
}

func (monitor *monitor) Write(chunk []byte) (int, error) {
	monitor.guard.Lock()
	monitor.buffer.Write(chunk)
	seen := strings.Contains(monitor.buffer.String(), monitor.marker)
	monitor.guard.Unlock()
	if seen {
		monitor.once.Do(func() { close(monitor.ready) })
	}
	return len(chunk), nil
}

func (monitor *monitor) report() string {
	monitor.guard.Lock()
	defer monitor.guard.Unlock()
	return monitor.buffer.String()
}
