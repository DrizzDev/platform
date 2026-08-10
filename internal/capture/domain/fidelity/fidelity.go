package fidelity

// Fidelity labels how truthfully a captured value reflects its source; an inferred value is never presented as exact.
type Fidelity string

const (
	Exact       Fidelity = "EXACT"
	Summary     Fidelity = "SUMMARY"
	Inferred    Fidelity = "INFERRED"
	Redacted    Fidelity = "REDACTED"
	Unavailable Fidelity = "UNAVAILABLE"
)

func (fidelity Fidelity) Valid() bool {
	switch fidelity {
	case Exact, Summary, Inferred, Unavailable, Redacted:
		return true
	default:
		return false
	}
}

func (fidelity Fidelity) String() string {
	return string(fidelity)
}
