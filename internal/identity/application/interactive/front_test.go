package interactive_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/interactive"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

type authorization struct {
	token  login.Token
	begin  error
	finish error
}

func (fake authorization) Begin(_ context.Context, _ login.Secret) (login.Redirect, error) {
	return login.Redirect{URL: "https://issuer.example/authorize"}, fake.begin
}

func (fake authorization) Finish(_ context.Context, _ login.Exchange) (login.Token, error) {
	return fake.token, fake.finish
}

type browser struct {
	prompt error
}

func (fake browser) Prompt(_ context.Context, _ login.Redirect) (login.Callback, error) {
	return login.Callback{Code: "code_123", State: "state_123"}, fake.prompt
}

type entropy struct {
	fail bool
}

func (fake entropy) Bytes(size int) ([]byte, error) {
	if fake.fail {
		return nil, errors.New("no entropy")
	}
	return bytes.Repeat([]byte{7}, size), nil
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) token() login.Token {
	fixture.test.Helper()
	account, failure := subject.New("google-oauth2|110000000000000000000")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return login.Token{
		Issuer:  "https://issuer.example/",
		Client:  "native",
		Subject: account,
		Method:  method.Browser,
		Refresh: []byte("refresh-token-bytes"),
		Issued:  time.Unix(1000, 0),
		Expiry:  time.Unix(2000, 0),
	}
}

func (fixture fixture) establish(options interactive.Options) (login.Token, error) {
	fixture.test.Helper()
	front, failure := interactive.New(options)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return front.Establish(context.Background())
}

func TestEstablish(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	token, failure := fixture.establish(interactive.Options{
		Authorization: authorization{token: fixture.token()},
		Browser:       browser{},
		Random:        entropy{},
	})
	if failure != nil {
		test.Fatalf("a valid establishment failed: %v", failure)
	}
	if token.Subject.String() != "google-oauth2|110000000000000000000" {
		test.Fatalf("subject = %q", token.Subject.String())
	}
}

func TestPropagates(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	cases := map[string]struct {
		options interactive.Options
		want    error
	}{
		"begin":  {interactive.Options{Authorization: authorization{begin: login.Unavailable{}}, Browser: browser{}, Random: entropy{}}, login.Unavailable{}},
		"prompt": {interactive.Options{Authorization: authorization{token: fixture.token()}, Browser: browser{prompt: login.Cancelled{}}, Random: entropy{}}, login.Cancelled{}},
		"finish": {interactive.Options{Authorization: authorization{finish: login.Rejected{}}, Browser: browser{}, Random: entropy{}}, login.Rejected{}},
	}
	for name, scenario := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			_, failure := fixture.establish(scenario.options)
			if !errors.Is(failure, scenario.want) {
				test.Fatalf("failure = %v, want %v", failure, scenario.want)
			}
		})
	}
}

func TestEntropy(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	_, failure := fixture.establish(interactive.Options{
		Authorization: authorization{token: fixture.token()},
		Browser:       browser{},
		Random:        entropy{fail: true},
	})
	if failure == nil {
		test.Fatal("an entropy failure was not surfaced")
	}
}
