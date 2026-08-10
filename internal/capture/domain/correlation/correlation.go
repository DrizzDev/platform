package correlation

import (
	"errors"

	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

// Correlation places one hop in the execution tree: a trace shared by the whole execution, this hop's span, its parent
// span (empty at the root), an ordering sequence, its cross-source bearings, and whether the parent link is inferred.
type Correlation struct {
	trace trace.Trace

	span   span.Span
	parent span.Span

	sequence int64
	mark     mark.Mark
	bearings bearings.Bearings
}

type Input struct {
	Trace trace.Trace

	Span   span.Span
	Parent span.Span

	Sequence int64
	Mark     mark.Mark
	Bearings bearings.Bearings
}

func New(input Input) (Correlation, error) {
	correlation := Correlation{
		span:     input.Span,
		mark:     input.Mark,
		trace:    input.Trace,
		parent:   input.Parent,
		bearings: input.Bearings,
		sequence: input.Sequence,
	}
	if failure := correlation.validate(); failure != nil {
		return Correlation{}, failure
	}
	return correlation, nil
}

func (correlation Correlation) Trace() trace.Trace {
	return correlation.trace
}

func (correlation Correlation) Span() span.Span {
	return correlation.span
}

func (correlation Correlation) Parent() span.Span {
	return correlation.parent
}

func (correlation Correlation) Bearings() bearings.Bearings {
	return correlation.bearings
}

func (correlation Correlation) Sequence() int64 {
	return correlation.sequence
}

func (correlation Correlation) Mark() mark.Mark {
	return correlation.mark
}

// Root reports the first hop of the execution — the one span with no parent.
func (correlation Correlation) Root() bool {
	return correlation.parent.Empty()
}

func (correlation Correlation) validate() error {
	switch {
	case correlation.trace.Empty():
		return errors.New("correlation trace is required")
	case correlation.span.Empty():
		return errors.New("correlation span is required")
	case correlation.sequence < 0:
		return errors.New("correlation sequence must not be negative")
	case !correlation.mark.Valid():
		return errors.New("correlation mark is invalid")
	case !correlation.parent.Empty() && correlation.parent == correlation.span:
		return errors.New("correlation parent must not be the span itself")
	default:
		return nil
	}
}
