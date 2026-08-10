package bridge

import "time"

// Options configures the channel. Digest is the pinned SHA-256 the binary must
// match before exec; Environment is the allowlisted set passed to the child on top
// of the device toolchain vars.
type Options struct {
	Location    string
	Digest      string
	Environment []string
	Timeout     time.Duration
}
