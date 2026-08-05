package organization_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/organization"
)

func TestOrganization(test *testing.T) {
	test.Parallel()

	value, failure := organization.New("organization_123")
	if failure != nil {
		test.Fatal(failure)
	}
	if value.String() != "organization_123" {
		test.Fatalf("organization = %q", value.String())
	}
	if _, failure := organization.New(""); failure == nil {
		test.Fatal("an invalid organization was accepted")
	}
}
