package journal

import (
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
	"github.com/DrizzDev/platform/internal/capture/domain/trace"
)

// Retained is a synchronized entry a reclaim pass may age out: which execution and step it is, its class (which fixes
// its retention), and when it was recorded. It carries no payload — reclaim decides on lifecycle, not content.
type Retained struct {
	Recorded time.Time
	Span     span.Span
	Trace    trace.Trace
	Category category.Category
}
