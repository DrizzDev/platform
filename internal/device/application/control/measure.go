package control

// Measure is a single numeric reading from a device, such as free disk space, or a code-only failure.
type Measure struct {
	outcome
	value int
}

func (measure Measure) Value() int {
	return measure.value
}
