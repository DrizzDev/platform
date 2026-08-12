package control

// Extent is a device's screen size in pixels, or a code-only failure.
type Extent struct {
	outcome
	width  int
	height int
}

func (extent Extent) Width() int {
	return extent.width
}

func (extent Extent) Height() int {
	return extent.height
}
