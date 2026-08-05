package reconcile

import fault "github.com/DrizzDev/platform/internal/identity/application/failure"

// Result is the outcome of a reconcile pass: how many candidates were reclaimed
// and how many were deferred, or a code-only failure when the backlog could not
// be read.
type Result struct {
	failure   *fault.Value
	reclaimed int
	deferred  int
}

func (result Result) Reclaimed() int {
	return result.reclaimed
}

func (result Result) Deferred() int {
	return result.deferred
}

func (result Result) Failed() bool {
	return result.failure != nil
}

func (result Result) Failure() (fault.Value, bool) {
	if result.failure == nil {
		return fault.Value{}, false
	}
	return *result.failure, true
}
