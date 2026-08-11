package catalog

// Entry describes one capability Drizz offers: its short internal name, the outward name an agent and a person see, a
// short human summary, and the inputs it accepts. It is the single description that both the command line and the agent
// connection advertise, so the two can never drift apart.
type Entry struct {
	name       string
	title      string
	summary    string
	parameters []Parameter
}

func (entry Entry) Name() string {
	return entry.name
}

// Title is the outward, unambiguous name the agent connection and the command line present.
func (entry Entry) Title() string {
	return entry.title
}

func (entry Entry) Summary() string {
	return entry.summary
}

func (entry Entry) Parameters() []Parameter {
	return append([]Parameter(nil), entry.parameters...)
}
