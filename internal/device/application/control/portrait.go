package control

import "github.com/DrizzDev/platform/internal/device/domain/capture"

// Portrait is a screen capture together with its on-screen element tree, or a code-only failure.
type Portrait struct {
	outcome
	capture   capture.Capture
	hierarchy string
}

func (portrait Portrait) Capture() capture.Capture {
	return portrait.capture
}

func (portrait Portrait) Hierarchy() string {
	return portrait.hierarchy
}
