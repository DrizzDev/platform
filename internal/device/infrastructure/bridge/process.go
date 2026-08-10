package bridge

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// passthrough is the device toolchain allowlist; no secret or unrelated host
// variable crosses into the child (SEC-008, SEC-013).
var passthrough = []string{
	"PATH", "HOME", "TMPDIR", "LANG",
	"ANDROID_HOME", "ANDROID_SDK_ROOT", "DEVELOPER_DIR",
}

type machine struct {
	command *exec.Cmd
	input   io.WriteCloser
	output  io.ReadCloser
	logs    io.ReadCloser
}

// spawn runs the sidecar with an explicit executable and arg list, an allowlisted
// environment, a fixed directory, and graceful-shutdown escalation — no shell (SEC-013).
func (options Options) spawn(scope context.Context) (machine, error) {

	command := exec.CommandContext(scope, options.Location)

	command.WaitDelay = grace
	command.Env = options.environment()
	command.Dir = filepath.Dir(options.Location)
	command.Cancel = func() error { return command.Process.Signal(os.Interrupt) }

	input, failure := command.StdinPipe()
	if failure != nil {
		return machine{}, Unavailable{}
	}
	output, failure := command.StdoutPipe()
	if failure != nil {
		return machine{}, Unavailable{}
	}
	logs, failure := command.StderrPipe()
	if failure != nil {
		return machine{}, Unavailable{}
	}
	if failure := command.Start(); failure != nil {
		return machine{}, Unavailable{}
	}
	return machine{command: command, input: input, output: output, logs: logs}, nil
}

func (options Options) environment() []string {
	values := make([]string, 0, len(passthrough)+len(options.Environment))
	for _, key := range passthrough {
		if value, found := os.LookupEnv(key); found {
			values = append(values, key+"="+value)
		}
	}
	return append(values, options.Environment...)
}
