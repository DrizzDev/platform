package credential_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

type builder struct {
	test *testing.T
}

func (builder builder) input() credential.Input {
	builder.test.Helper()
	account, failure := subject.New("google-oauth2|110000000000")
	if failure != nil {
		builder.test.Fatal(failure)
	}
	handle, failure := session.New("session_123")
	if failure != nil {
		builder.test.Fatal(failure)
	}
	return credential.Input{
		Issuer:   "https://issuer.example/",
		Client:   "native",
		Handle:   "handle_1234567890",
		Subject:  account,
		Session:  handle,
		Method:   method.Browser,
		Refresh:  []byte("refresh-token-bytes"),
		Issued:   time.Unix(1000, 0),
		Expiry:   time.Unix(2000, 0),
		Revision: 1,
		Schema:   1,
	}
}

func TestRecord(test *testing.T) {
	test.Parallel()

	input := builder{test: test}.input()
	record, failure := credential.New(input)
	if failure != nil {
		test.Fatal(failure)
	}
	if record.Issuer() != input.Issuer || record.Method() != method.Browser {
		test.Fatalf("record = %+v", record)
	}
	if string(record.Key()) != input.Session.String()+"#1#"+input.Handle {
		test.Fatalf("key = %q", record.Key())
	}

	bytes := record.Refresh()
	bytes[0] = 0
	if record.Refresh()[0] == 0 {
		test.Fatal("refresh bytes are not defensively copied")
	}
}

func TestInvalid(test *testing.T) {
	test.Parallel()

	base := builder{test: test}.input()
	cases := map[string]func(credential.Input) credential.Input{
		"issuer":   func(input credential.Input) credential.Input { input.Issuer = ""; return input },
		"client":   func(input credential.Input) credential.Input { input.Client = ""; return input },
		"handle":   func(input credential.Input) credential.Input { input.Handle = ""; return input },
		"subject":  func(input credential.Input) credential.Input { input.Subject = subject.Subject{}; return input },
		"session":  func(input credential.Input) credential.Input { input.Session = session.Session{}; return input },
		"method":   func(input credential.Input) credential.Input { input.Method = method.Method("OTHER"); return input },
		"refresh":  func(input credential.Input) credential.Input { input.Refresh = nil; return input },
		"oversize": func(input credential.Input) credential.Input { input.Refresh = make([]byte, 3000); return input },
		"expiry":   func(input credential.Input) credential.Input { input.Expiry = time.Unix(500, 0); return input },
		"revision": func(input credential.Input) credential.Input { input.Revision = 0; return input },
		"schema":   func(input credential.Input) credential.Input { input.Schema = 0; return input },
	}
	for name, mutate := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			if _, failure := credential.New(mutate(base)); failure == nil {
				test.Fatal("invalid record was accepted")
			}
		})
	}
}
