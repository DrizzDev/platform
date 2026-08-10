package origin

// Origin marks whether a captured fact is an authoritative Drizz execution event or an agent-side host observation.
type Origin string

const (
	Host       Origin = "HOST"
	Capability Origin = "CAPABILITY"
)

func (origin Origin) Valid() bool {
	switch origin {
	case Capability, Host:
		return true
	default:
		return false
	}
}

func (origin Origin) String() string {
	return string(origin)
}
