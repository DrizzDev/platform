package control

import "github.com/DrizzDev/platform/internal/device/domain/capture"

// Observation is one captured screen, or a code-only failure.
type Observation struct {
	outcome
	capture capture.Capture
}

func (observation Observation) Capture() capture.Capture {
	return observation.capture
}
