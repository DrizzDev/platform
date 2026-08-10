package janitor

import (
	"context"
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
	"github.com/DrizzDev/platform/internal/capture/domain/span"
)

// sifted separates reclaim candidates: aged records are past their retention and removed every pass; fresh records are
// still within retention and removed only under disk pressure. A protected (leased) trace falls in neither.
type sifted struct {
	aged  []span.Span
	fresh []span.Span
}

// protected is the set of traces whose lease has not yet expired, so a reclaim pass leaves a running execution alone.
func (janitor Janitor) protected(scope context.Context) (map[string]bool, error) {
	claims, failure := janitor.archive.Leases(scope)

	if failure != nil {
		return nil, failure
	}
	now := janitor.clock.Now()
	live := map[string]bool{}
	for _, claim := range claims {
		if claim.Until.After(now) {
			live[claim.Trace.String()] = true
		}
	}
	return live, nil
}

// sift gathers the synchronized candidates and the active leases, then splits the candidates into records past their
// retention and records still within it, leaving every leased trace out of both.
func (janitor Janitor) sift(scope context.Context) (sifted, error) {
	protected, failure := janitor.protected(scope)

	if failure != nil {
		return sifted{}, failure
	}
	settled, failure := janitor.archive.Settled(scope)
	if failure != nil {
		return sifted{}, failure
	}
	now := janitor.clock.Now()
	var parts sifted
	for _, candidate := range settled {
		if protected[candidate.Trace.String()] {
			continue
		}
		policy, failure := janitor.catalogue.Policy(candidate.Category)
		if failure != nil {
			continue
		}
		if now.Sub(candidate.Recorded) >= policy.Retention() {
			parts.aged = append(parts.aged, candidate.Span)
		} else {
			parts.fresh = append(parts.fresh, candidate.Span)
		}
	}
	return parts, nil
}

// sweep deletes every stored object no remaining entry references, sparing any still within the grace window so a
// just-written object racing its own entry is never reaped.
func (janitor Janitor) sweep(scope context.Context) error {
	referenced, failure := janitor.archive.Digests(scope)

	if failure != nil {
		return failure
	}
	live := map[string]bool{}
	for _, reference := range referenced {
		live[reference.String()] = true
	}
	now := janitor.clock.Now()
	keep := func(reference digest.Digest, modified time.Time) bool {
		return live[reference.String()] || now.Sub(modified) < janitor.grace
	}
	return janitor.vault.Prune(scope, keep)
}
