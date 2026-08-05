package browser

import (
	"context"
	"os/exec"
	"runtime"
)

var _ Opener = System{}

// System opens the URL through the platform browser launcher with an explicit
// executable and a single argument, never a shell (SEC-013).
type System struct{}

func (system System) Open(scope context.Context, target string) error {
	return system.command(scope, target).Start()
}

func (System) command(scope context.Context, target string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.CommandContext(scope, "rundll32", "url.dll,FileProtocolHandler", target)
	case "darwin":
		return exec.CommandContext(scope, "open", target)
	default:
		return exec.CommandContext(scope, "xdg-open", target)
	}
}
