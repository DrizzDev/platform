package lease_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/coordination/lease"
)

func TestLease(test *testing.T) {
	test.Parallel()

	held, failure := lease.New(lease.Input{Owner: "process", Expiry: time.Unix(2000, 0)})
	if failure != nil {
		test.Fatal(failure)
	}
	if held.Owner() != "process" {
		test.Fatalf("owner = %q", held.Owner())
	}
	if !held.Held(time.Unix(1000, 0)) {
		test.Fatal("lease was not held before its expiry")
	}
	if held.Held(time.Unix(3000, 0)) {
		test.Fatal("lease was held after its expiry")
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	cases := map[string]lease.Input{
		"owner":  {Expiry: time.Unix(2000, 0)},
		"expiry": {Owner: "process"},
	}
	for name, input := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			if _, failure := lease.New(input); failure == nil {
				test.Fatal("invalid lease was accepted")
			}
		})
	}
}
