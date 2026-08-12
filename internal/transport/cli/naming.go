package cli

import (
	"strings"
	"unicode"
)

// slug turns an outward capability name into its command-line form — words separated by dashes, lowercased — so the
// command line and the agent connection name the same capability consistently from the one catalogue. TakeScreenshot
// becomes take-screenshot.
func (Command) slug(title string) string {
	runes := []rune(title)
	var builder strings.Builder
	for index, letter := range runes {
		if index > 0 && unicode.IsUpper(letter) {
			following := index+1 < len(runes) && unicode.IsLower(runes[index+1])
			if unicode.IsLower(runes[index-1]) || following {
				builder.WriteRune('-')
			}
		}
		builder.WriteRune(unicode.ToLower(letter))
	}
	return builder.String()
}
