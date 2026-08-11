package cli

import (
	"strings"
	"unicode"
)

// slug turns an outward capability name into its command-line form — words separated by dashes, lowercased — so the
// command line and the agent connection name the same capability consistently from the one catalogue. TakeScreenshot
// becomes take-screenshot.
func (Command) slug(title string) string {
	var builder strings.Builder
	for index, letter := range title {
		if index > 0 && unicode.IsUpper(letter) {
			builder.WriteRune('-')
		}
		builder.WriteRune(unicode.ToLower(letter))
	}
	return builder.String()
}
