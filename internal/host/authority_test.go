package host_test

import (
	"context"
	"testing"

	"github.com/DrizzDev/platform/internal/host"
	"github.com/DrizzDev/platform/internal/identity/application/login"
)

func TestAuthorityPermits(test *testing.T) {
	test.Parallel()

	made, failure := host.Authority("")
	if failure != nil {
		test.Fatal(failure)
	}
	tenant, failure := made.Authorize(context.Background(), login.Grant{})
	if failure != nil {
		test.Fatalf("a no-cloud authority denied the sign-in: %v", failure)
	}
	if tenant.Name != "" {
		test.Fatalf("a no-cloud authority surfaced a tenant: %q", tenant.Name)
	}
}
