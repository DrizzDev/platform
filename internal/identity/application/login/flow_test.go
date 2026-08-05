package login_test

import (
	"context"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/standing"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

type establishment struct {
	token login.Token
	fail  error
}

func (fake establishment) Establish(context.Context) (login.Token, error) {
	return fake.token, fake.fail
}

type publication struct {
	receipt   login.Receipt
	publish   error
	retracted bool
}

func (fake *publication) Publish(_ context.Context, _ login.Candidate) (login.Receipt, error) {
	return fake.receipt, fake.publish
}

func (fake *publication) Retract(context.Context, time.Time) error {
	fake.retracted = true
	return nil
}

type authority struct {
	tenant login.Tenant
	fail   error
}

func (fake authority) Authorize(context.Context, login.Grant) (login.Tenant, error) {
	return fake.tenant, fake.fail
}

type clock struct{}

func (clock) Now() time.Time {
	return time.Unix(1500, 0)
}

type setup struct {
	establishment login.Establishment
	publication   login.Publication
	authority     login.Authority
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

func (fixture fixture) receipt() login.Receipt {
	fixture.test.Helper()
	account, failure := subject.New("google-oauth2|110000000000000000000")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	handle, failure := session.New("session_1234567890")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return login.Receipt{Subject: account, Session: handle, Method: method.Browser, Expiry: time.Unix(2000, 0)}
}

func (fixture fixture) base() setup {
	return setup{
		establishment: establishment{token: fixture.token()},
		publication:   &publication{receipt: fixture.receipt()},
		authority:     authority{tenant: login.Tenant{Name: "Acme"}},
	}
}

func (fixture fixture) flow(setup setup) login.Flow {
	fixture.test.Helper()
	made, failure := login.New(login.Options{
		Establishment: setup.establishment,
		Publication:   setup.publication,
		Authority:     setup.authority,
		Clock:         clock{},
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestRun(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	result := fixture.flow(fixture.base()).Run(context.Background(), login.Input{})
	if result.Failed() {
		test.Fatalf("a valid login failed: %+v", result)
	}
	if result.Subject().String() != "google-oauth2|110000000000000000000" {
		test.Fatalf("subject = %q", result.Subject().String())
	}
	if result.Session().String() != "session_1234567890" {
		test.Fatalf("session = %q", result.Session().String())
	}
	if result.Method() != method.Browser || result.Standing() != standing.Active {
		test.Fatalf("method = %q, standing = %q", result.Method(), result.Standing())
	}
	if !result.Expiry().Equal(time.Unix(2000, 0)) {
		test.Fatalf("expiry = %v", result.Expiry())
	}
	if result.Organization() != "Acme" {
		test.Fatalf("organization = %q", result.Organization())
	}
}

func TestForbiddenRetracts(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	store := &publication{receipt: fixture.receipt()}
	result := fixture.flow(setup{
		establishment: establishment{token: fixture.token()},
		publication:   store,
		authority:     authority{fail: login.Forbidden{}},
	}).Run(context.Background(), login.Input{})
	value, denied := result.Failure()
	if !denied || value.Code() != code.Forbidden {
		test.Fatalf("failure = %+v %v, want forbidden", value, denied)
	}
	if !store.retracted {
		test.Fatal("a forbidden login did not retract the published credential")
	}
}

func TestDenied(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	cases := map[string]struct {
		mutate func(setup) setup
		want   code.Code
	}{
		"cancelled":   {func(base setup) setup { base.establishment = establishment{fail: login.Cancelled{}}; return base }, code.Cancelled},
		"rejected":    {func(base setup) setup { base.establishment = establishment{fail: login.Rejected{}}; return base }, code.Rejected},
		"unavailable": {func(base setup) setup { base.establishment = establishment{fail: login.Unavailable{}}; return base }, code.Unavailable},
		"conflict":    {func(base setup) setup { base.publication = &publication{publish: login.Conflict{}}; return base }, code.Conflict},
		"storage":     {func(base setup) setup { base.publication = &publication{publish: login.Storage{}}; return base }, code.Storage},
		"forbidden":   {func(base setup) setup { base.authority = authority{fail: login.Forbidden{}}; return base }, code.Forbidden},
		"context":     {func(base setup) setup { base.establishment = establishment{fail: context.Canceled}; return base }, code.Cancelled},
	}
	for name, scenario := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			result := fixture.flow(scenario.mutate(fixture.base())).Run(context.Background(), login.Input{})
			value, denied := result.Failure()
			if !denied || value.Code() != scenario.want {
				test.Fatalf("failure = %+v %v, want %q", value, denied, scenario.want)
			}
		})
	}
}
