package identifier

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const limit = 256

// Identifier is a validated correlation id; the zero value is the explicit absence of an id (an unset parent), so
// callers test Empty rather than comparing strings.
type Identifier struct {
	value string
}

func New(value string) (Identifier, error) {
	identifier := Identifier{value: value}
	if failure := identifier.validate(); failure != nil {
		return Identifier{}, failure
	}
	return identifier, nil
}

func (identifier Identifier) String() string {
	return identifier.value
}

func (identifier Identifier) Empty() bool {
	return identifier.value == ""
}

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
