package journal

import (
	"errors"

	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/correlation"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
)

// Entry is one durable record in the execution journal: where it sits in the trace, what it is, how truthful it is,
// its classified payload, and its synchronization state.
type Entry struct {
	state       State
	payload     []byte
	origin      origin.Origin
	fidelity    fidelity.Fidelity
	category    category.Category
	correlation correlation.Correlation
}

type Input struct {
	State       State
	Payload     []byte
	Origin      origin.Origin
	Fidelity    fidelity.Fidelity
	Category    category.Category
	Correlation correlation.Correlation
}

func New(input Input) (Entry, error) {
	entry := Entry{
		state:       input.State,
		origin:      input.Origin,
		fidelity:    input.Fidelity,
		category:    input.Category,
		correlation: input.Correlation,
		payload:     append([]byte(nil), input.Payload...),
	}
	if failure := entry.validate(); failure != nil {
		return Entry{}, failure
	}
	return entry, nil
}

func (entry Entry) Correlation() correlation.Correlation {
	return entry.correlation
}

func (entry Entry) Origin() origin.Origin {
	return entry.origin
}

func (entry Entry) Fidelity() fidelity.Fidelity {
	return entry.fidelity
}

func (entry Entry) Category() category.Category {
	return entry.category
}

func (entry Entry) Payload() []byte {
	return append([]byte(nil), entry.payload...)
}

func (entry Entry) State() State {
	return entry.state
}

func (entry Entry) validate() error {
	switch {
	case entry.correlation.Trace().Empty():
		return errors.New("entry correlation is required")
	case !entry.origin.Valid():
		return errors.New("entry origin is invalid")
	case !entry.fidelity.Valid():
		return errors.New("entry fidelity is invalid")
	case !entry.category.Valid():
		return errors.New("entry category is invalid")
	case !entry.state.Valid():
		return errors.New("entry state is invalid")
	default:
		return nil
	}
}
