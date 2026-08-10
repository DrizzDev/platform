package bridge

// Fault is a device-bridge application error surfaced only by kind; the raw sidecar message is never carried here (SEC-008).
type Fault struct {
	kind string
}

func (Fault) Error() string {
	return "the device operation failed"
}

func (fault Fault) Kind() string {
	return fault.kind
}

// The transport errors are typed structs so callers match them without string comparison and without an errFoo
// sentinel the name gate would reject.

type Closed struct{}

func (Closed) Error() string {
	return "device channel is closed"
}

type Unavailable struct{}

func (Unavailable) Error() string {
	return "device sidecar is unavailable"
}

type Mismatch struct{}

func (Mismatch) Error() string {
	return "device sidecar protocol is unsupported"
}

type Oversize struct{}

func (Oversize) Error() string {
	return "device sidecar sent an oversized frame"
}
