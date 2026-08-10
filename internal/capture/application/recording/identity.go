package recording

import (
	"crypto/rand"
	"encoding/hex"
)

// mint returns a random 128-bit hex id for a trace or span. Uniqueness comes from randomness;
// ordering is carried by the entry sequence, so the id need not be time-sortable.
func (Recorder) mint() (string, error) {
	buffer := make([]byte, 16)
	if _, failure := rand.Read(buffer); failure != nil {
		return "", failure
	}
	return hex.EncodeToString(buffer), nil
}
