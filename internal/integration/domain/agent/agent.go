package agent

// Kind is the stable, one-word-per-token name of a supported agent application. It is the lookup key every surface
// uses, so the command line, the catalog, and a record all refer to the same agent.
type Kind string

func (kind Kind) String() string {
	return string(kind)
}

// Agent describes one supported agent application as data: which agent it is, the outward title a person sees, where
// its configuration file lives (an anchor plus path segments, resolved later by an adapter), the file's format, whether
// its stdio entry carries an explicit type, and whether the agent must be restarted to notice a change. Everything the
// installer needs to wire an agent is here, so adding an agent is a new row rather than new code.
type Agent struct {
	kind       Kind
	title      string
	base       Base
	segments   []string
	collection string
	dialect    Dialect
	hooking    Hooking
	commanding Commanding
	typed      bool
	restart    bool
}

func (agent Agent) Kind() Kind {
	return agent.kind
}

// Title is the outward, human-readable name of the agent application.
func (agent Agent) Title() string {
	return agent.title
}

func (agent Agent) Base() Base {
	return agent.base
}

// Segments are the path elements, under the resolved base, that locate the agent's configuration file.
func (agent Agent) Segments() []string {
	return append([]string(nil), agent.segments...)
}

// Collection is the key under which the agent's configuration groups its MCP servers — "mcpServers" for the
// JSON agents, "mcp_servers" for Codex. Naming it in the descriptor keeps the merge engine format-agnostic.
func (agent Agent) Collection() string {
	return agent.collection
}

func (agent Agent) Dialect() Dialect {
	return agent.dialect
}

// Hooking describes how Drizz registers for this agent's turn events. Its Supported reports whether the agent has a
// hook mechanism at all.
func (agent Agent) Hooking() Hooking {
	return agent.hooking
}

// Commanding describes where this agent keeps its user-invocable commands, so the installer can place the Drizz
// command file. An agent with no command surface leaves it unset.
func (agent Agent) Commanding() Commanding {
	return agent.commanding
}

// Typed reports whether the agent's stdio server entry must carry an explicit "type" of "stdio". Claude Code requires
// it; others default to stdio and omit it.
func (agent Agent) Typed() bool {
	return agent.typed
}

// Restart reports whether the agent must be fully restarted before it notices a configuration change. A desktop
// application that reads its configuration once at launch does; a command line that reads it each session does not.
func (agent Agent) Restart() bool {
	return agent.restart
}
