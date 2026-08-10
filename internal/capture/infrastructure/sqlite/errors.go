package sqlite

// Corrupt reports that a stored journal row could not be reconstructed into a valid value and was isolated rather than
// returned as a partial result.
type Corrupt struct{}

func (Corrupt) Error() string {
	return "a journal row is corrupt"
}
