package janitor

import (
	"context"

	"github.com/DrizzDev/platform/internal/capture/domain/span"
)

// Run performs one reclaim pass. It ages out synchronized records past their retention and deletes the artifacts they
// leave unreferenced. If the store is still over budget it reclaims synchronized records within retention as well, and
// when only unacknowledged evidence remains it degrades new writes rather than deleting it.
func (janitor Janitor) Run(scope context.Context) (Report, error) {
	parts, failure := janitor.sift(scope)

	if failure != nil {
		return Report{}, failure
	}

	if failure := janitor.reclaim(scope, parts.aged); failure != nil {
		return Report{}, failure
	}
	report := Report{Reclaimed: len(parts.aged)}

	used, failure := janitor.vault.Footprint(scope)
	if failure != nil {
		return Report{}, failure
	}
	if used > janitor.ceiling {
		if failure := janitor.reclaim(scope, parts.fresh); failure != nil {
			return Report{}, failure
		}
		report.Reclaimed += len(parts.fresh)
		if used, failure = janitor.vault.Footprint(scope); failure != nil {
			return Report{}, failure
		}
	}
	report.Degraded, failure = janitor.regulate(scope, used)
	return report, failure
}

// reclaim deletes the selected entries, then sweeps the artifacts no remaining entry references. Both steps are
// idempotent, so a pass interrupted between them safely resumes on the next run.
func (janitor Janitor) reclaim(scope context.Context, steps []span.Span) error {
	if failure := janitor.archive.Discard(scope, steps); failure != nil {
		return failure
	}
	return janitor.sweep(scope)
}

// regulate holds the write gate to the budget with hysteresis: over the ceiling it degrades writes and alerts; at or
// below the relief mark it restores them; between the two it leaves the gate as it stands.
func (janitor Janitor) regulate(scope context.Context, used int64) (bool, error) {
	switch {
	case used > janitor.ceiling:
		if failure := janitor.vault.Restrict(scope); failure != nil {
			return false, failure
		}
		return true, janitor.notifier.Alert(scope, Pressure{Used: used, Ceiling: janitor.ceiling})
	case used <= janitor.relief:
		return false, janitor.vault.Admit(scope)
	default:
		return false, nil
	}
}
