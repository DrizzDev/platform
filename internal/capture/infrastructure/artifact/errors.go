package artifact

// Oversize reports that an artifact exceeded the configured size ceiling and was refused.
type Oversize struct{}

func (Oversize) Error() string {
	return "the artifact exceeds the size ceiling"
}

// Integrity reports that a stored artifact's bytes no longer match its digest.
type Integrity struct{}

func (Integrity) Error() string {
	return "the artifact failed integrity verification"
}

// Absent reports that no artifact is stored for the requested digest.
type Absent struct{}

func (Absent) Error() string {
	return "no artifact is stored for the digest"
}
