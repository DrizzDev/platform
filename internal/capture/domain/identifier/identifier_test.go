package identifier_test

import (
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/identifier"
)

func TestIdentifier(test *testing.T) {
	test.Parallel()

	value, failure := identifier.New("01HEXEC0001")
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != "01HEXEC0001" || value.Empty() {
		test.Fatalf("identifier = %q empty=%v", value.String(), value.Empty())
	}
	if !(identifier.Identifier{}).Empty() {
		test.Fatal("the zero identifier is not reported empty")
	}
}

func TestIdentifierRejects(test *testing.T) {
	test.Parallel()

	rejected := map[string]string{"empty": "", "control": "a\x00b", "long": strings.Repeat("a", 257)}
	for name, value := range rejected {
		if _, failure := identifier.New(value); failure == nil {
			test.Fatalf("%s identifier was accepted", name)
		}
	}
}
