package lobby

import (
	"context"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
)

// Observe holds one host observation in the pending window and drops any that have expired, so the window stays
// bounded in time. It is called from the host hook path, before any Drizz capability call.
func (lobby Lobby) Observe(scope context.Context, observation journal.Observation) error {
	if failure := lobby.register.Admit(scope, observation); failure != nil {
		return failure
	}
	return lobby.register.Expire(scope, lobby.cutoff())
}
