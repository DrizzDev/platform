package connect

import "github.com/DrizzDev/platform/internal/integration/domain/agent"

// State is the outcome of wiring or unwiring one agent. It is a small, stable vocabulary a surface can render without
// interpreting a cause.
type State string

const (
	// Connected means Drizz was newly registered with the agent.
	Connected State = "CONNECTED"
	// Updated means Drizz was already registered and its entry was refreshed.
	Updated State = "UPDATED"
	// Removed means the Drizz entry was taken out of the agent's configuration.
	Removed State = "REMOVED"
	// Missing means the agent application is not installed, so there was nothing to do.
	Missing State = "MISSING"
	// Ready means the agent is installed but Drizz is not connected to it yet.
	Ready State = "READY"
	// Captured means Drizz is registered for the agent's turn events.
	Captured State = "CAPTURED"
	// Cleared means Drizz's turn-event registration was removed from the agent.
	Cleared State = "CLEARED"
	// Incapable means the agent is installed but has no hook mechanism, so context capture is not available for it.
	Incapable State = "INCAPABLE"
	// Failed means the agent is installed but the change could not be completed.
	Failed State = "FAILED"
)

func (state State) String() string {
	return string(state)
}

// Selection chooses which agents an operation applies to: one named agent, or all supported agents.
type Selection struct {
	Kind agent.Kind
	All  bool
}

// Outcome is the result for one agent: which agent, its human title, what happened, whether Drizz is capturing that
// agent's turn events, whether the agent must be restarted to notice the change, and a short safe detail for a failure.
type Outcome struct {
	kind      agent.Kind
	title     string
	state     State
	detail    string
	restart   bool
	capturing bool
}

func (outcome Outcome) Kind() agent.Kind {
	return outcome.kind
}

func (outcome Outcome) Title() string {
	return outcome.title
}

func (outcome Outcome) State() State {
	return outcome.state
}

func (outcome Outcome) Restart() bool {
	return outcome.restart
}

// Capturing reports whether Drizz is registered for this agent's turn events.
func (outcome Outcome) Capturing() bool {
	return outcome.capturing
}

func (outcome Outcome) Detail() string {
	return outcome.detail
}

// Report is the per-agent result of one installer operation.
type Report struct {
	outcomes []Outcome
}

func (report Report) Outcomes() []Outcome {
	return append([]Outcome(nil), report.outcomes...)
}
