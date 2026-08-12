package agent

// catalogue is the canonical, ordered set of supported agent applications. It is read-only: List copies before
// returning it. Adding support for another agent — Gemini's peers, a new editor's assistant — is a new row here plus,
// only if it introduces a new file format, one reader and writer for that dialect.
var catalogue = []Agent{
	{
		kind:       "claude-code",
		title:      "Claude Code",
		base:       Home,
		segments:   []string{".claude.json"},
		collection: "mcpServers",
		dialect:    Json,
		typed:      true,
		restart:    false,
		hooking: Hooking{
			base:     Home,
			segments: []string{".claude", "settings.json"},
			dialect:  Json,
			style:    Claude,
			channel:  Stdin,
			slots:    []Slot{Prompt, Turn},
		},
	},
	{
		kind:       "claude-desktop",
		title:      "Claude Desktop",
		base:       Config,
		segments:   []string{"Claude", "claude_desktop_config.json"},
		collection: "mcpServers",
		dialect:    Json,
		typed:      false,
		restart:    true,
		hooking:    Hooking{style: None},
	},
	{
		kind:       "codex",
		title:      "Codex",
		base:       Home,
		segments:   []string{".codex", "config.toml"},
		collection: "mcp_servers",
		dialect:    Toml,
		typed:      false,
		restart:    false,
		hooking: Hooking{
			base:     Home,
			segments: []string{".codex", "config.toml"},
			dialect:  Toml,
			style:    Codex,
			channel:  Argv,
			slots:    []Slot{Turn},
		},
	},
	{
		kind:       "gemini",
		title:      "Gemini",
		base:       Home,
		segments:   []string{".gemini", "settings.json"},
		collection: "mcpServers",
		dialect:    Json,
		typed:      false,
		restart:    false,
		hooking: Hooking{
			base:     Home,
			segments: []string{".gemini", "settings.json"},
			dialect:  Json,
			style:    Gemini,
			channel:  Stdin,
			slots:    []Slot{Prompt, Turn},
		},
	},
}

// Catalog is the single, ordered list of supported agents. The command line reads it to enumerate agents, and the
// installer reads each agent's data to wire it, so the two never disagree about what is supported or how.
type Catalog struct {
	agents []Agent
}

func New() Catalog {
	return Catalog{agents: catalogue}
}

func (catalog Catalog) List() []Agent {
	return append([]Agent(nil), catalog.agents...)
}

func (catalog Catalog) Lookup(kind Kind) (Agent, bool) {
	for _, agent := range catalog.agents {
		if agent.kind == kind {
			return agent, true
		}
	}
	return Agent{}, false
}
