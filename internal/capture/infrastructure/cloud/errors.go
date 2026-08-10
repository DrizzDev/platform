package cloud

import "strconv"

// transient reports an upload the cloud could not accept this time — a server fault or an unexpected status — that may
// succeed on a later attempt. It carries only the numeric status, never a response body, so no cloud detail leaks.
type transient struct {
	status int
}

func (transient transient) Error() string {
	return "the upload failed transiently with status " + strconv.Itoa(transient.status)
}

// unreachable reports that the request never got a response — a dial, connection, or read fault below HTTP. It is
// retryable and carries no detail, so a peer address or transport message never leaks to a caller.
type unreachable struct{}

func (unreachable) Error() string {
	return "the upload could not reach the cloud"
}
