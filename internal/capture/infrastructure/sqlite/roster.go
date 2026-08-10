package sqlite

import (
	"time"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/bearings"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/fidelity"
	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
	"github.com/DrizzDev/platform/internal/capture/domain/origin"
)

// slot is one pending-window row.
type slot struct {
	ordinal     int64
	moment      int64
	reference   int64
	session     string
	turn        string
	call        string
	connection  string
	execution   string
	capability  string
	fingerprint string
	process     string
	origin      string
	fidelity    string
	category    string
	payload     []byte
}

// observations is the ordered layout of the pending table.
var observations = layout[journal.Observation, slot]{
	{
		name:  "session",
		read:  func(slot *slot) any { return &slot.session },
		write: func(observation journal.Observation) any { return observation.Bearings.Session().String() },
	},
	{
		name:  "turn",
		read:  func(slot *slot) any { return &slot.turn },
		write: func(observation journal.Observation) any { return observation.Bearings.Turn().String() },
	},
	{
		name:  "call",
		read:  func(slot *slot) any { return &slot.call },
		write: func(observation journal.Observation) any { return observation.Bearings.Call().String() },
	},
	{
		name:  "connection",
		read:  func(slot *slot) any { return &slot.connection },
		write: func(observation journal.Observation) any { return observation.Bearings.Connection().String() },
	},
	{
		name:  "execution",
		read:  func(slot *slot) any { return &slot.execution },
		write: func(observation journal.Observation) any { return observation.Bearings.Execution().String() },
	},
	{
		name:  "capability",
		read:  func(slot *slot) any { return &slot.capability },
		write: func(observation journal.Observation) any { return observation.Bearings.Capability().String() },
	},
	{
		name:  "fingerprint",
		read:  func(slot *slot) any { return &slot.fingerprint },
		write: func(observation journal.Observation) any { return observation.Fingerprint.String() },
	},
	{
		name:  "process",
		read:  func(slot *slot) any { return &slot.process },
		write: func(observation journal.Observation) any { return observation.Process.String() },
	},
	{
		name:  "ordinal",
		read:  func(slot *slot) any { return &slot.ordinal },
		write: func(observation journal.Observation) any { return observation.Ordinal },
	},
	{
		name:  "moment",
		read:  func(slot *slot) any { return &slot.moment },
		write: func(observation journal.Observation) any { return observation.Moment.Unix() },
	},
	{
		name:  "origin",
		read:  func(slot *slot) any { return &slot.origin },
		write: func(observation journal.Observation) any { return observation.Origin.String() },
	},
	{
		name:  "fidelity",
		read:  func(slot *slot) any { return &slot.fidelity },
		write: func(observation journal.Observation) any { return observation.Fidelity.String() },
	},
	{
		name:  "category",
		read:  func(slot *slot) any { return &slot.category },
		write: func(observation journal.Observation) any { return observation.Category.String() },
	},
	{
		name:  "payload",
		read:  func(slot *slot) any { return &slot.payload },
		write: func(observation journal.Observation) any { return observation.Payload },
	},
}

// held reconstructs the stored observation, validating each identity and fingerprint; a value that cannot be rebuilt
// is a corruption of throwaway data and is reported so the caller can drop it.
func (slot *slot) held() (journal.Held, error) {

	identities, failure := slot.identities()
	if failure != nil {
		return journal.Held{}, Corrupt{}
	}

	process, failure := slot.optional(slot.process)
	if failure != nil {
		return journal.Held{}, Corrupt{}
	}

	fingerprint := digest.Digest{}
	if slot.fingerprint != "" {
		if fingerprint, failure = digest.New(slot.fingerprint); failure != nil {
			return journal.Held{}, Corrupt{}
		}
	}
	return journal.Held{Reference: slot.reference, Observation: journal.Observation{
		Process:     process,
		Bearings:    identities,
		Fingerprint: fingerprint,
		Ordinal:     slot.ordinal,
		Payload:     slot.payload,
		Moment:      time.Unix(slot.moment, 0),
		Origin:      origin.Origin(slot.origin),
		Fidelity:    fidelity.Fidelity(slot.fidelity),
		Category:    category.Category(slot.category),
	}}, nil
}

func (slot *slot) identities() (bearings.Bearings, error) {
	raw := []string{slot.session, slot.turn, slot.call, slot.connection, slot.execution, slot.capability}

	ids := make([]identifier.Identifier, len(raw))
	for index, value := range raw {
		made, failure := slot.optional(value)
		if failure != nil {
			return bearings.Bearings{}, failure
		}
		ids[index] = made
	}
	return bearings.New(bearings.Input{
		Session: ids[0], Turn: ids[1], Call: ids[2], Connection: ids[3], Execution: ids[4], Capability: ids[5],
	}), nil
}

func (*slot) optional(value string) (identifier.Identifier, error) {
	if value == "" {
		return identifier.Identifier{}, nil
	}
	return identifier.New(value)
}
