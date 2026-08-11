package system

import "time"

// Clock reads the host wall clock for the recording lease deadline; it is injected so recording stays deterministic
// under test.
type Clock struct{}

func New() Clock {
	return Clock{}
}

func (Clock) Now() time.Time {
	return time.Now()
}
