package organization_test

import (
	"context"
	"testing"

	"github.com/DrizzDev/platform/internal/identity/application/organization"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
	tenant "github.com/DrizzDev/platform/internal/identity/domain/organization"
)

type resolver struct {
	org  tenant.Organization
	fail error
}

func (fake resolver) Resolve(context.Context) (tenant.Organization, error) {
	return fake.org, fake.fail
}

type fixture struct {
	test *testing.T
}

func (fixture fixture) org() tenant.Organization {
	fixture.test.Helper()
	made, failure := tenant.New("organization_123")
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func (fixture fixture) flow(resolve organization.Resolver) organization.Flow {
	fixture.test.Helper()
	made, failure := organization.New(organization.Options{Resolver: resolve})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestResolve(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	result := fixture.flow(resolver{org: fixture.org()}).Resolve(context.Background(), organization.Input{})
	if result.Failed() {
		value, _ := result.Failure()
		test.Fatalf("resolution failed: %s", value.Code())
	}
	if result.Organization().String() != "organization_123" {
		test.Fatalf("organization = %q", result.Organization().String())
	}
}

func TestDenied(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	cases := map[string]struct {
		fail error
		want code.Code
	}{
		"forbidden":   {organization.Forbidden{}, code.Forbidden},
		"required":    {organization.Required{}, code.Required},
		"unavailable": {organization.Unavailable{}, code.Unavailable},
		"context":     {context.Canceled, code.Cancelled},
	}
	for name, scenario := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			result := fixture.flow(resolver{fail: scenario.fail}).Resolve(context.Background(), organization.Input{})
			value, denied := result.Failure()
			if !denied || value.Code() != scenario.want {
				test.Fatalf("failure = %+v %v, want %q", value, denied, scenario.want)
			}
		})
	}
}
