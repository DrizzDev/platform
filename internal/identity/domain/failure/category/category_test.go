package category_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/failure/category"
)

func TestValid(test *testing.T) {
	test.Parallel()

	values := []category.Category{
		category.Authentication,
		category.Authorization,
		category.Account,
		category.Storage,
		category.Logout,
		category.Internal,
	}
	for _, value := range values {
		if !value.Valid() {
			test.Fatalf("category %q was rejected", value)
		}
		if value.String() != string(value) {
			test.Fatalf("category = %q", value.String())
		}
	}
}

func TestUnknown(test *testing.T) {
	test.Parallel()

	if category.Category("OTHER").Valid() {
		test.Fatal("an unknown category was accepted")
	}
}
