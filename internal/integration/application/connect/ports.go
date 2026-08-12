package connect

import (
	"context"

	"github.com/DrizzDev/platform/internal/capture/application/recording"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
	"github.com/DrizzDev/platform/internal/integration/domain/server"
)

// Task is one agent to wire and the server entry to write into it, bundled so the store's methods stay within the
// one-typed-input shape the rest of the codebase uses. It is exported because a store adapter in another package
// implements the port and must name the input.
type Task struct {
	Agent  agent.Agent
	Server server.Server
}

// resolver yields the path of the running Drizz executable, so the entry written into an agent's configuration
// launches exactly this installed binary. An infrastructure adapter satisfies it.
type resolver interface {
	Locate() (string, error)
}

// store reads and writes one agent application's configuration file. Detect reports whether the agent is present on
// this machine; Wired reports whether Drizz is already registered; Connect merges the Drizz entry, preserving every
// other setting, and confirms the write; Disconnect removes only the Drizz entry. Per-dialect adapters satisfy it, so
// the flow never knows whether a file is JSON or TOML.
type store interface {
	Detect(agent.Agent) (bool, error)
	Wired(agent.Agent) (bool, error)
	Connect(context.Context, Task) error
	Disconnect(context.Context, agent.Agent) error
	Captures(agent.Agent) (bool, error)
	Capture(context.Context, Task) error
	Uncapture(context.Context, agent.Agent) error
}

// recorder opens one execution record; the capture recorder satisfies it. Recording is observational, so a failure to
// open or write a record never fails the installer action it observes.
type recorder interface {
	Begin() (*recording.Execution, error)
}
