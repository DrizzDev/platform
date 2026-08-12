package agent

// Style names how an agent's hook registration is shaped in its configuration file. Each style is a different layout —
// Claude and Gemini keep a hooks object keyed by event, Codex keeps a single notify program — so a writer is selected
// by style rather than by branching on the agent everywhere.
type Style string

const (
	// None means the agent has no hook mechanism to register.
	None Style = "NONE"
	// Claude is the Claude Code hooks object: events like UserPromptSubmit and Stop, each holding command handlers.
	Claude Style = "CLAUDE"
	// Gemini is the Gemini CLI hooks object: events like BeforeAgent and AfterAgent, each holding command handlers.
	Gemini Style = "GEMINI"
	// Codex is the Codex notify program: a single command invoked when a turn completes.
	Codex Style = "CODEX"
)

func (style Style) Valid() bool {
	switch style {
	case None, Claude, Gemini, Codex:
		return true
	default:
		return false
	}
}

func (style Style) String() string {
	return string(style)
}
