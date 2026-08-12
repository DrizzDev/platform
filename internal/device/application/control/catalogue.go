package control

// Catalogue is a set of named items read from the device layer, such as emulator image names, or a code-only failure.
type Catalogue struct {
	outcome
	names []string
}

func (catalogue Catalogue) Names() []string {
	return append([]string(nil), catalogue.names...)
}
