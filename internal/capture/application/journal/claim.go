package journal

import (
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

// Claim is an active execution's lease on its trace: a reclaim pass skips a trace whose claim has not yet expired, so
// a running execution's records are protected until it ends or the lease lapses.
type Claim struct {
	Until time.Time
	Trace trace.Trace
}
