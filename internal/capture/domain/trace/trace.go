package trace

import "github.com/DrizzDev/platform/internal/capture/domain/identifier"

// Trace is the root id shared by an entire execution — the anchor that groups every hop first call to last.
type Trace struct {
	identifier.Identifier
}

func New(value string) (Trace, error) {
	inner, failure := identifier.New(value)
	if failure != nil {
		return Trace{}, failure
	}
	return Trace{Identifier: inner}, nil
}
