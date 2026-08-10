package lobby

import (
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
)

// Call is a Drizz capability call to match against the pending window: the same cross-source attributes an observation
// carries, so the two correlate exactly by a shared identifier or inferred by proximity.
type Call struct {
	Ordinal     int64
	Moment      time.Time
	Fingerprint digest.Digest
	Bearings    bearings.Bearings
	Process     identifier.Identifier
}
