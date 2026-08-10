package journal

// State is where a recorded entry sits in the synchronization lifecycle.
type State string

const (
	Synced  State = "SYNCED"
	Pending State = "PENDING"
)

func (state State) Valid() bool {
	switch state {
	case Pending, Synced:
		return true
	default:
		return false
	}
}

func (state State) String() string {
	return string(state)
}
