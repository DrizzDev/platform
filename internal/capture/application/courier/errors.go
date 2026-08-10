package courier

// Rejected reports an upload the cloud refused for a reason that will not change on retry — a validation,
// authorization, or integrity failure. The courier surfaces it instead of retrying.
type Rejected struct{}

func (Rejected) Error() string {
	return "the upload was rejected and must not be retried"
}

func (Rejected) terminal() bool {
	return true
}

// terminal is implemented by a failure that must not be retried.
type terminal interface {
	terminal() bool
}
