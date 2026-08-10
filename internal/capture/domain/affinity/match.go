package affinity

import "github.com/DrizzDev/platform/internal/capture/domain/mark"

// Match decides how this call relates to an earlier observation. A shared stable identifier is an exact link.
// Otherwise, an observation that preceded the call within the call's window and agrees on normalized input or process
// is an inferred link — never relabelled exact. Anything else is no match.
func (call Signal) Match(candidate Signal) (mark.Mark, bool) {
	if call.bearings.Shares(candidate.bearings) {
		return mark.Exact, true
	}
	if call.infers(candidate) {
		return mark.Inferred, true
	}
	return "", false
}

func (call Signal) infers(candidate Signal) bool {
	if candidate.ordinal > call.ordinal || !call.near(candidate) {
		return false
	}
	return call.fingerprint.Same(candidate.fingerprint) || call.process.Same(candidate.process)
}

func (call Signal) near(candidate Signal) bool {
	gap := call.moment.Sub(candidate.moment)
	if gap < 0 {
		gap = -gap
	}
	return gap <= call.window
}
