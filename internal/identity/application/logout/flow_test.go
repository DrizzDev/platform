package logout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/credential"
	"github.com/DrizzDev/platform/internal/identity/application/logout"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
	"github.com/DrizzDev/platform/internal/identity/domain/method"
	"github.com/DrizzDev/platform/internal/identity/domain/session"
	"github.com/DrizzDev/platform/internal/identity/domain/subject"
)

type vault struct {
	record credential.Record
	active error
}

func (fake vault) Active(_ context.Context) (credential.Record, error) {
	return fake.record, fake.active
}

type publication struct {
	retract error
}

func (fake publication) Retract(_ context.Context, _ time.Time) error {
	return fake.retract
}

type revocation struct {
	revoke error
}

func (fake revocation) Revoke(_ context.Context, _ credential.Record) error {
	return fake.revoke
}

type clock struct{}

func (clock) Now() time.Time {
	return time.Unix(1500, 0)
}

type setup struct {
	vault       vault
	publication publication
	revocation  revocation
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) record() credential.Record {
	fixture.test.Helper()
	account, failure := subject.New("google-oauth2|first")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	handle, failure := session.New("LOCAL")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	record, failure := credential.New(credential.Input{
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
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return record
}

func (fixture fixture) base() setup {
	return setup{vault: vault{record: fixture.record()}}
}

func (fixture fixture) flow(setup setup) logout.Flow {
	fixture.test.Helper()
	made, failure := logout.New(logout.Options{
		Vault:       setup.vault,
		Publication: setup.publication,
		Revocation:  setup.revocation,
		Clock:       clock{},
	})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestRun(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	if result := fixture.flow(fixture.base()).Run(context.Background(), logout.Input{}); result.Failed() {
		test.Fatalf("a clean logout failed: %+v", result)
	}
}

func TestIdempotent(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	setup := setup{vault: vault{active: logout.Missing{}}}
	if result := fixture.flow(setup).Run(context.Background(), logout.Input{}); result.Failed() {
		test.Fatalf("logging out when signed out should succeed: %+v", result)
	}
}

func TestOutcomes(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	cases := map[string]struct {
		mutate func(setup) setup
		want   code.Code
	}{
		"storage":   {func(base setup) setup { base.vault = vault{active: logout.Storage{}}; return base }, code.Storage},
		"retract":   {func(base setup) setup { base.publication = publication{retract: errors.New("locked")}; return base }, code.Failed},
		"partial":   {func(base setup) setup { base.revocation = revocation{revoke: errors.New("unreachable")}; return base }, code.Partial},
		"cancelled": {func(base setup) setup { base.vault = vault{active: context.Canceled}; return base }, code.Cancelled},
	}
	for name, scenario := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			result := fixture.flow(scenario.mutate(fixture.base())).Run(context.Background(), logout.Input{})
			value, failed := result.Failure()
			if !failed || value.Code() != scenario.want {
				test.Fatalf("failure = %+v %v, want %q", value, failed, scenario.want)
			}
		})
	}
}
