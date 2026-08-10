package affinity

import (
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
)

// Signal is the matchable trace of one observation or call: its cross-source identities, a fingerprint of its
// normalized input, the process that produced it, when it was seen, its position in the observed order, and — for a
// call — the window within which it will accept an inferred match.
type Signal struct {
	ordinal     int64
	moment      time.Time
	window      time.Duration
	fingerprint digest.Digest
	bearings    bearings.Bearings
	process     identifier.Identifier
}

type Input struct {
	Ordinal     int64
	Moment      time.Time
	Window      time.Duration
	Fingerprint digest.Digest
	Bearings    bearings.Bearings
	Process     identifier.Identifier
}

func New(input Input) Signal {
	return Signal{
		moment:      input.Moment,
		window:      input.Window,
		process:     input.Process,
		ordinal:     input.Ordinal,
		bearings:    input.Bearings,
		fingerprint: input.Fingerprint,
	}
}
