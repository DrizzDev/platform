package mcp

import (
	"context"
	"errors"
	"io"
)

// outcome classifies the result of a serve session fail-safe: a session never
// classifies an external error as an internal defect, so no serve outcome is
// reported to the error sink.
func (execution execution) outcome(failure error) outcome {
	switch {
	case failure == nil:
		return success
	case errors.Is(failure, context.Canceled) || errors.Is(failure, io.EOF):
		return cancelled
	case execution.rejected(failure):
		return rejected
	default:
		return interrupted
	}
}

func (execution execution) rejected(failure error) bool {
	var big excess
	if errors.As(failure, &big) {
		return true
	}
	var bad malformed
	return errors.As(failure, &bad)
}
