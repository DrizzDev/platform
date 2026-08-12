package intake

import "github.com/DrizzDev/platform/internal/integration/domain/agent"

// Event is one inbound hook notification: which agent sent it, which turn moment it marks, and the text the agent
// exposed for that moment — a person's prompt, or the agent's own final message. The text is whatever the agent chose
// to reveal; Drizz records it as an observed host fact, never as its own authoritative output.
type Event struct {
	Agent agent.Kind
	Slot  agent.Slot
	Text  string
}
