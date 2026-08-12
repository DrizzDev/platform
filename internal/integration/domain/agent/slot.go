package agent

// Slot is a moment in an agent's turn that Drizz can be notified about. Prompt is when a person submits a prompt;
// Turn is when the agent finishes responding. An agent supports the slots its hook system exposes — some expose both,
// some only one.
type Slot string

const (
	Prompt Slot = "prompt"
	Turn   Slot = "turn"
)

func (slot Slot) String() string {
	return string(slot)
}
