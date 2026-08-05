package cleanup

type Reason string

const (
	Logout     Reason = "LOGOUT"
	Rejected   Reason = "REJECTED"
	Superseded Reason = "SUPERSEDED"
)

func (reason Reason) Valid() bool {
	switch reason {
	case Rejected, Superseded, Logout:
		return true
	default:
		return false
	}
}

func (reason Reason) String() string {
	return string(reason)
}
