package cleanup

type State string

const (
	Pending State = "PENDING"
	Blocked State = "BLOCKED"
)

func (state State) Valid() bool {
	switch state {
	case Pending, Blocked:
		return true
	default:
		return false
	}
}

func (state State) String() string {
	return string(state)
}
