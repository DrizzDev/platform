package session_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/session"
)

func TestSession(test *testing.T) {
	test.Parallel()

	value, failure := session.New("session_123")
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != "session_123" {
		test.Fatalf("session = %q", value.String())
	}
	if _, failure := session.New(""); failure == nil {
		test.Fatal("an invalid session was accepted")
	}
}
