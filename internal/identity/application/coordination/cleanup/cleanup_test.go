package cleanup_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/cleanup"
)

type builder struct{}

func (builder builder) input() cleanup.Input {
	return cleanup.Input{
		Key:      "session_123#1",
		Reason:   cleanup.Superseded,
		State:    cleanup.Pending,
		Attempts: 0,
		Next:     time.Unix(1000, 0),
		Deadline: time.Unix(2000, 0),
		Created:  time.Unix(500, 0),
	}
}

func TestRecord(test *testing.T) {
	test.Parallel()

	record, failure := cleanup.New(builder{}.input())
	if failure != nil {
		test.Fatal(failure)
	}
	if record.Reason() != cleanup.Superseded || record.Blocked() {
		test.Fatalf("record = %+v", record)
	}

	input := builder{}.input()
	input.State = cleanup.Blocked
	blocked, failure := cleanup.New(input)
	if failure != nil {
		test.Fatal(failure)
	}
	if !blocked.Blocked() {
		test.Fatal("a blocked record did not report blocked")
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	cases := map[string]func(cleanup.Input) cleanup.Input{
		"key":      func(input cleanup.Input) cleanup.Input { input.Key = ""; return input },
		"reason":   func(input cleanup.Input) cleanup.Input { input.Reason = cleanup.Reason("OTHER"); return input },
		"state":    func(input cleanup.Input) cleanup.Input { input.State = cleanup.State("OTHER"); return input },
		"created":  func(input cleanup.Input) cleanup.Input { input.Created = time.Time{}; return input },
		"deadline": func(input cleanup.Input) cleanup.Input { input.Deadline = time.Time{}; return input },
	}
	for name, mutate := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			if _, failure := cleanup.New(mutate(builder{}.input())); failure == nil {
				test.Fatal("invalid cleanup record was accepted")
			}
		})
	}
}

func TestReason(test *testing.T) {
	test.Parallel()

	for _, value := range []cleanup.Reason{cleanup.Rejected, cleanup.Superseded, cleanup.Logout} {
		if !value.Valid() {
			test.Fatalf("reason %q was rejected", value)
		}
	}
	if cleanup.Reason("OTHER").Valid() {
		test.Fatal("an unknown reason was accepted")
	}
}

func TestState(test *testing.T) {
	test.Parallel()

	for _, value := range []cleanup.State{cleanup.Pending, cleanup.Blocked} {
		if !value.Valid() {
			test.Fatalf("state %q was rejected", value)
		}
	}
	if cleanup.State("OTHER").Valid() {
		test.Fatal("an unknown state was accepted")
	}
}
