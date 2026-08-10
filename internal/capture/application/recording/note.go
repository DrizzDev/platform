package recording

import (
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
)

// Note is one thing to record: what it is, how truthful it is, a small typed payload,
// and optionally the bytes of a large artifact to store alongside it.
type Note struct {
	Payload  []byte
	Artifact []byte
	Origin   origin.Origin
	Fidelity fidelity.Fidelity
	Category category.Category
}
