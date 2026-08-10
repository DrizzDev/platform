package catalog

// Entry describes one capability Drizz offers: its name, a short human summary, and the inputs it accepts. It is the
// single description that both the command line and the agent connection advertise.
type Entry struct {
	name       string
	summary    string
	parameters []Parameter
}

func (entry Entry) Name() string {
	return entry.name
}

func (entry Entry) Summary() string {
	return entry.summary
}

func (entry Entry) Parameters() []Parameter {
	return append([]Parameter(nil), entry.parameters...)
}
