package mark

// Mark records whether a correlation link is directly observed or inferred,
// so an inferred relationship is never. presented as exact.
type Mark string

const (
	Exact    Mark = "EXACT"
	Inferred Mark = "INFERRED"
)

func (mark Mark) Valid() bool {
	switch mark {
	case Exact, Inferred:
		return true
	default:
		return false
	}
}

func (mark Mark) String() string {
	return string(mark)
}
