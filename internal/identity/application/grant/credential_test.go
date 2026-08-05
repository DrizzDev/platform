package grant_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/grant"
)

func TestCredential(test *testing.T) {
	test.Parallel()

	credential, failure := grant.New(grant.Input{
		Token:  []byte("access-token"),
		Expiry: time.Unix(2000, 0),
	})
	if failure != nil {
		test.Fatal(failure)
	}
	if credential.Expired(time.Unix(1000, 0)) {
		test.Fatal("credential reported expired before its expiry")
	}
	if !credential.Expired(time.Unix(3000, 0)) {
		test.Fatal("credential did not report expired after its expiry")
	}

	bytes := credential.Token()
	bytes[0] = 0
	if credential.Token()[0] == 0 {
		test.Fatal("token is not defensively copied")
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	cases := map[string]grant.Input{
		"token":  {Expiry: time.Unix(2000, 0)},
		"expiry": {Token: []byte("access-token")},
	}
	for name, input := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			if _, failure := grant.New(input); failure == nil {
				test.Fatal("invalid grant was accepted")
			}
		})
	}
}
