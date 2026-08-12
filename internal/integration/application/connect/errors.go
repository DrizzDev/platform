package connect

// The port errors below are the installer-owned failure contract. A store adapter returns them so the flow can
// classify an outcome into a stable state without inspecting an operating-system cause.

// Absent means the agent application is not present on this machine, so there is nothing to configure.
type Absent struct{}

func (Absent) Error() string {
	return "the agent application is not installed"
}

func (Absent) state() State {
	return Missing
}

// Malformed means the agent's existing configuration file could not be parsed. The installer refuses to overwrite a
// file it cannot understand, so no other setting is ever lost.
type Malformed struct{}

func (Malformed) Error() string {
	return "the agent configuration file could not be read"
}

func (Malformed) state() State {
	return Failed
}

// Locked means the configuration file could not be read or written, typically a permission or ownership problem.
type Locked struct{}

func (Locked) Error() string {
	return "the agent configuration file could not be written"
}

func (Locked) state() State {
	return Failed
}

// Unsupported means the agent has no hook mechanism Drizz can register for, so context capture is not available for it.
type Unsupported struct{}

func (Unsupported) Error() string {
	return "the agent has no hook mechanism"
}

func (Unsupported) state() State {
	return Incapable
}

// Occupied means the agent already routes its single turn-notification to a different program, so Drizz will not
// overwrite it. It applies to agents, like Codex, that allow only one notification program.
type Occupied struct{}

func (Occupied) Error() string {
	return "the agent already has a different notification program configured"
}

func (Occupied) state() State {
	return Failed
}
