package lobby

import "time"

// Options configures a Lobby. Window is how long an unmatched observation lives and the horizon for an inferred time
// match; Capacity bounds how many observations one match considers. Each defaults in one place when left zero.
type Options struct {
	Capacity int
	Clock    clock
	Register register
	Window   time.Duration
}

func (options Options) span() time.Duration {
	if options.Window == 0 {
		return window
	}
	return options.Window
}

func (options Options) size() int {
	if options.Capacity == 0 {
		return capacity
	}
	return options.Capacity
}
