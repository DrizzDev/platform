package action_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/failure/action"
)

func TestValid(test *testing.T) {
	test.Parallel()

	for _, value := range []action.Action{action.Signin, action.Retry, action.Logout, action.None} {
		if !value.Valid() {
			test.Fatalf("action %q was rejected", value)
		}
		if value.String() != string(value) {
			test.Fatalf("action = %q", value.String())
		}
	}
}

func TestUnknown(test *testing.T) {
	test.Parallel()

	if action.Action("OTHER").Valid() {
		test.Fatal("an unknown action was accepted")
	}
}
