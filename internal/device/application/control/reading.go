package control

// Reading is one text read from a device, such as the on-screen element tree or the foreground app, or a code-only
// failure.
type Reading struct {
	outcome
	text string
}

func (reading Reading) Text() string {
	return reading.text
}
