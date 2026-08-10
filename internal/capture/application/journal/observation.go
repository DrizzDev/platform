package journal

import (
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
)

// Observation is one host-side observation held in the bounded pending window until a Drizz capability call claims it
// or it expires. It carries the attributes that match it to a call and the note to record once it is claimed. It is
// never synchronized.
type Observation struct {
	Ordinal     int64
	Payload     []byte
	Moment      time.Time
	Fingerprint digest.Digest
	Origin      origin.Origin
	Fidelity    fidelity.Fidelity
	Category    category.Category
	Bearings    bearings.Bearings
	Process     identifier.Identifier
}

// Held is a stored observation and the reference the store uses to evict it once it is claimed or expired.
type Held struct {
	Reference   int64
	Observation Observation
}

// Window bounds a read of the pending store: only observations at or after the cutoff,
// and never more than the capacity, so the matchable set stays bounded.
type Window struct {
	Capacity int
	Cutoff   time.Time
}
