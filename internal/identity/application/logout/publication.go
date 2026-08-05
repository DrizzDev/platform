package logout

import (
	"context"
	"time"
)

// Publication clears the active credential locally in one atomic step: it drops
// the active pointer, advances the epoch, and queues the key for deletion.
type Publication interface {
	Retract(context.Context, time.Time) error
}
