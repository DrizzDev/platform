package digest_test

import (
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

func TestDigest(test *testing.T) {
	test.Parallel()

	valid := strings.Repeat("a", 64)
	value, failure := digest.New(valid)
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != valid || value.Empty() {
		test.Fatalf("digest = %q empty=%v", value.String(), value.Empty())
	}
	if !(digest.Digest{}).Empty() {
		test.Fatal("the zero digest is not reported empty")
	}
}

func TestDigestRejects(test *testing.T) {
	test.Parallel()

	rejected := map[string]string{
		"short":  "abc",
		"long":   strings.Repeat("a", 65),
		"upper":  strings.Repeat("A", 64),
		"nonhex": strings.Repeat("g", 64),
	}
	for name, value := range rejected {
		if _, failure := digest.New(value); failure == nil {
			test.Fatalf("%s digest was accepted", name)
		}
	}
}
