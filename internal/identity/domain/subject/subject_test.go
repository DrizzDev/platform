package subject_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

func TestSubject(test *testing.T) {
	test.Parallel()

	value, failure := subject.New("google-oauth2|110000000000")
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != "google-oauth2|110000000000" {
		test.Fatalf("subject = %q", value.String())
	}
	if _, failure := subject.New(""); failure == nil {
		test.Fatal("an invalid subject was accepted")
	}
}
