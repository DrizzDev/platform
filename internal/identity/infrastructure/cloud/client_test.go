package cloud_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/identity/application/grant"
	"github.com/DrizzDev/platform/internal/identity/application/login"
	"github.com/DrizzDev/platform/internal/identity/application/organization"
	"github.com/DrizzDev/platform/internal/identity/application/session"
	"github.com/DrizzDev/platform/internal/identity/infrastructure/cloud"
)

type provider struct {
	credential grant.Credential
	fail       error
}

func (fake provider) Access(context.Context, session.Input) (grant.Credential, error) {
	return fake.credential, fake.fail
}

type clock struct{}

func (clock) Now() time.Time { return time.Unix(1000, 0) }

type fixture struct {
	test *testing.T
}

func (fixture fixture) grant() grant.Credential {
	fixture.test.Helper()
	credential, failure := grant.New(grant.Input{Token: []byte("access-token"), Expiry: time.Unix(9000, 0)})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return credential
}

func (fixture fixture) client(options cloud.Options) cloud.Client {
	fixture.test.Helper()
	options.Agent = &http.Client{}
	options.Clock = clock{}
	made, failure := cloud.New(options)
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestResolve(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer access-token" {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"id":42,"name":"Acme"}`))
	}))
	defer server.Close()

	organization, failure := fixture.client(cloud.Options{Provider: provider{credential: fixture.grant()}, Base: server.URL}).Resolve(context.Background())
	if failure != nil {
		test.Fatalf("resolution failed: %v", failure)
	}
	if organization.String() != "42" {
		test.Fatalf("organization = %q", organization.String())
	}
}

func TestOutcomes(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	cases := map[string]struct {
		status int
		want   error
	}{
		"forbidden":   {http.StatusForbidden, organization.Forbidden{}},
		"required":    {http.StatusUnauthorized, organization.Required{}},
		"unavailable": {http.StatusBadGateway, organization.Unavailable{}},
	}
	for name, scenario := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				writer.WriteHeader(scenario.status)
			}))
			defer server.Close()
			_, failure := fixture.client(cloud.Options{Provider: provider{credential: fixture.grant()}, Base: server.URL}).Resolve(context.Background())
			if !errors.Is(failure, scenario.want) {
				test.Fatalf("failure = %v, want %v", failure, scenario.want)
			}
		})
	}
}

func (fixture fixture) gate(base string) cloud.Client {
	fixture.test.Helper()
	made, failure := cloud.New(cloud.Options{Agent: &http.Client{}, Base: base})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return made
}

func TestAuthorize(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer access-token" {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = writer.Write([]byte(`{"id":42,"name":"Acme"}`))
	}))
	defer server.Close()

	tenant, failure := fixture.gate(server.URL).Authorize(context.Background(),
		login.Grant{Token: []byte("access-token"), Expiry: time.Unix(9000, 0)})
	if failure != nil {
		test.Fatalf("authorization failed: %v", failure)
	}
	if tenant.Name != "Acme" {
		test.Fatalf("tenant = %q", tenant.Name)
	}
}

func TestForbids(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	_, failure := fixture.gate(server.URL).Authorize(context.Background(),
		login.Grant{Token: []byte("access-token"), Expiry: time.Unix(9000, 0)})
	if !errors.Is(failure, login.Forbidden{}) {
		test.Fatalf("failure = %v, want Forbidden", failure)
	}
}

func TestReauth(test *testing.T) {
	test.Parallel()

	fixture := fixture{test: test}
	_, failure := fixture.client(cloud.Options{Provider: provider{fail: session.Missing{}}, Base: "http://cloud.invalid"}).Resolve(context.Background())
	if !errors.Is(failure, organization.Required{}) {
		test.Fatalf("failure = %v, want Required", failure)
	}
}
