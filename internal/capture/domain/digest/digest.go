package digest

import "errors"

// Digest is a content-addressed SHA-256 identity: 64 lowercase hex characters.
// The zero value is the explicit absence of a digest.
type Digest struct {
	value string
}

func New(value string) (Digest, error) {
	digest := Digest{value: value}
	if failure := digest.validate(); failure != nil {
		return Digest{}, failure
	}
	return digest, nil
}

func (digest Digest) String() string {
	return digest.value
}

func (digest Digest) Empty() bool {
	return digest.value == ""
}

// Same reports whether both digests are present and equal, so a shared content fingerprint can infer a match.
func (digest Digest) Same(other Digest) bool {
	return !digest.Empty() && digest == other
}

func (digest Digest) validate() error {
	if len(digest.value) != 64 {
		return errors.New("digest must be 64 hex characters")
	}
	for _, letter := range digest.value {
		digit := letter >= '0' && letter <= '9'
		lower := letter >= 'a' && letter <= 'f'
		if !digit && !lower {
			return errors.New("digest must be lowercase hex")
		}
	}
	return nil
}
