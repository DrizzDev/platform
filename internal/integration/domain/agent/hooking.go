package agent

// Hooking describes how Drizz registers for an agent's turn events: where the agent keeps its hook configuration
// (an anchor plus path segments, which may differ from where it keeps its MCP servers), the file's format, the layout
// style of the hook block, how the agent delivers an event to a hook program, and which turn moments it can report.
// An agent with no hook system has the None style and no slots.
type Hooking struct {
	base     Base
	segments []string
	dialect  Dialect
	style    Style
	channel  Channel
	slots    []Slot
}

// Supported reports whether the agent has a hook mechanism Drizz can register for.
func (hooking Hooking) Supported() bool {
	return hooking.style != None && hooking.style != ""
}

func (hooking Hooking) Base() Base {
	return hooking.base
}

func (hooking Hooking) Segments() []string {
	return append([]string(nil), hooking.segments...)
}

func (hooking Hooking) Dialect() Dialect {
	return hooking.dialect
}

func (hooking Hooking) Style() Style {
	return hooking.style
}

func (hooking Hooking) Channel() Channel {
	return hooking.channel
}

// Slots are the turn moments this agent can report — a prompt submission, a completed turn, or both.
func (hooking Hooking) Slots() []Slot {
	return append([]Slot(nil), hooking.slots...)
}
