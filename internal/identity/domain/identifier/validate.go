package identifier

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const limit = 256

func (identifier Identifier) validate() error {
	switch {
	case identifier.value == "":
		return errors.New("identifier is required")
	case len(identifier.value) > limit:
		return errors.New("identifier exceeds the maximum length")
	case !utf8.ValidString(identifier.value):
		return errors.New("identifier is not valid UTF-8")
	case strings.ContainsFunc(identifier.value, unicode.IsControl):
		return errors.New("identifier contains a control character")
	default:
		return nil
	}
}
