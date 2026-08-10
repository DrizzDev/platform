package lobby

import (
	"context"

	"github.com/DrizzDev/platform/internal/capture/application/journal"
	"github.com/DrizzDev/platform/internal/capture/domain/affinity"
)

// Activate matches a Drizz capability call against the pending window and returns the observations it claims, each
// tagged exact or inferred. Claimed observations are evicted so they are activated at most once; stale ones expire.
func (lobby Lobby) Activate(scope context.Context, call Call) ([]Claimed, error) {
	held, failure := lobby.register.Waiting(scope, journal.Window{Cutoff: lobby.cutoff(), Capacity: lobby.capacity})
	if failure != nil {
		return nil, failure
	}
	seeker := affinity.New(affinity.Input{
		Bearings:    call.Bearings,
		Fingerprint: call.Fingerprint,
		Process:     call.Process,
		Moment:      call.Moment,
		Ordinal:     call.Ordinal,
		Window:      lobby.window,
	})
	var claimed []Claimed
	var references []int64
	for _, candidate := range held {
		link, matched := seeker.Match(lobby.signal(candidate.Observation))
		if !matched {
			continue
		}
		claimed = append(claimed, Claimed{Observation: candidate.Observation, Mark: link})
		references = append(references, candidate.Reference)
	}
	if failure := lobby.register.Evict(scope, references); failure != nil {
		return nil, failure
	}
	if failure := lobby.register.Expire(scope, lobby.cutoff()); failure != nil {
		return nil, failure
	}
	return claimed, nil
}

func (Lobby) signal(observation journal.Observation) affinity.Signal {
	return affinity.New(affinity.Input{
		Moment:      observation.Moment,
		Ordinal:     observation.Ordinal,
		Process:     observation.Process,
		Bearings:    observation.Bearings,
		Fingerprint: observation.Fingerprint,
	})
}
