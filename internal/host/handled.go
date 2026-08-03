package host

// handled marks an error whose diagnostics the runtime already emitted, so the
// process boundary maps it to an exit code without reporting it again.
type handled struct {
	cause error
}

func (handled handled) Error() string {
	return handled.cause.Error()
}

func (handled handled) Unwrap() error {
	return handled.cause
}
