package agent

// Commanding describes where an agent keeps its user-invocable slash commands: the anchor and the path segments that
// locate the Drizz command file, so the installer can place a `/author` command without the domain resolving a real
// directory. An agent with no command surface leaves it unset, and Supported reports false.
type Commanding struct {
	base     Base
	segments []string
}

// Supported reports whether the agent exposes a command surface Drizz can install into.
func (commanding Commanding) Supported() bool {
	return len(commanding.segments) > 0
}

func (commanding Commanding) Base() Base {
	return commanding.base
}

func (commanding Commanding) Segments() []string {
	return append([]string(nil), commanding.segments...)
}
