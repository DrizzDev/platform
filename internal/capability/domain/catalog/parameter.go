package catalog

// Parameter describes one input a capability accepts. Every surface advertises the same name and meaning, so the
// command line and the agent connection can never describe the same input differently.
type Parameter struct {
	name    string
	summary string
}

func (parameter Parameter) Name() string {
	return parameter.name
}

func (parameter Parameter) Summary() string {
	return parameter.summary
}
