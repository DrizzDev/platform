package identifier_test

import (
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/identifier"
)

func TestIdentifier(test *testing.T) {
	test.Parallel()

	value, failure := identifier.New("google-oauth2|110000000000")
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != "google-oauth2|110000000000" {
		test.Fatalf("identifier = %q", value.String())
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	cases := map[string]string{
		"empty":    "",
		"oversize": strings.Repeat("a", 257),
		"control":  "bad\x00value",
		"encoding": "\xff\xfe",
	}
	for name, value := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			if _, failure := identifier.New(value); failure == nil {
				test.Fatal("invalid identifier was accepted")
			}
		})
	}
}
