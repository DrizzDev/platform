package lobby

import (
	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/mark"
)

// Claimed is one pending observation a capability call has taken, tagged with how it was matched — exact when the two
// shared an identifier, inferred when matched by proximity. The caller records it into the activated execution.
type Claimed struct {
	Mark        mark.Mark
	Observation journal.Observation
}
