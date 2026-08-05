package logout

import fault "github.com/DrizzDev/platform/internal/identity/application/failure"

// Result is the outcome of a logout. A cleared local session is a success; a
// code-only failure carries a partial or unavailable outcome.
type Result struct {
	failure *fault.Value
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
