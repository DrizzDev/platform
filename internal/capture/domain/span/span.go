package span

import "github.com/DrizzDev/platform/internal/capture/domain/identifier"

// Span is the id of one hop in the execution tree; the zero value marks the root, which has no parent.
type Span struct {
	identifier.Identifier
}

func New(value string) (Span, error) {
	inner, failure := identifier.New(value)
	if failure != nil {
		return Span{}, failure
	}
	return Span{Identifier: inner}, nil
}
