package carrier

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

// asset is the carried device helper: its bytes and the hex SHA-256 they must hash to before the helper is written or
// run. Verification is the trust boundary — an unverified helper is never executed.
type asset struct {
	bytes  []byte
	digest string
}

// verify reports whether a candidate's SHA-256 equals the pinned digest, compared in constant time so a mismatch
// cannot be probed by timing.
func (asset asset) verify(candidate []byte) bool {
	expected, failure := hex.DecodeString(strings.TrimSpace(asset.digest))
	if failure != nil {
		return false
	}
	sum := sha256.Sum256(candidate)
	return subtle.ConstantTimeCompare(sum[:], expected) == 1
}
