package lobby

import "time"

// Lobby is the bounded pending window where host observations wait to be claimed by a Drizz capability call. It holds
// only what a match may need: observations expire after the window and never leave it unless a call activates them.
type Lobby struct {
	capacity int
	clock    clock
	register register
	window   time.Duration
}

func (lobby Lobby) cutoff() time.Time {
	return lobby.clock.Now().Add(-lobby.window)
}
