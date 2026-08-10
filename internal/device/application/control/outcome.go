package control

import "github.com/DrizzDev/platform/internal/device/domain/code"

// outcome is the failure carried by every result. An empty code means success.
type outcome struct {
	failure code.Code
}

func (outcome outcome) Failed() bool {
	return outcome.failure != ""
}

func (outcome outcome) Failure() (code.Code, bool) {
	if outcome.failure == "" {
		return "", false
	}
	return outcome.failure, true
}
