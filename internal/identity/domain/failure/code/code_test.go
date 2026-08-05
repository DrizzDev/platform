package code_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/identity/domain/failure/action"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/category"
	"github.com/DrizzDev/platform/internal/identity/domain/failure/code"
)

type expectation struct {
	code      code.Code
	category  category.Category
	action    action.Action
	retryable bool
}

func TestMapping(test *testing.T) {
	test.Parallel()

	cases := []expectation{
		{code.Required, category.Authentication, action.Signin, false},
		{code.Cancelled, category.Authentication, action.None, false},
		{code.Forbidden, category.Authorization, action.None, false},
		{code.Unavailable, category.Authentication, action.Retry, true},
		{code.Rejected, category.Authentication, action.Retry, true},
		{code.Conflict, category.Account, action.Logout, false},
		{code.Storage, category.Storage, action.None, false},
		{code.Partial, category.Logout, action.None, false},
		{code.Failed, category.Internal, action.Retry, true},
	}
	for _, item := range cases {
		if !item.code.Valid() {
			test.Fatalf("code %q was rejected", item.code)
		}
		if item.code.Category() != item.category {
			test.Fatalf("%q category = %q", item.code, item.code.Category())
		}
		if item.code.Action() != item.action {
			test.Fatalf("%q action = %q", item.code, item.code.Action())
		}
		if item.code.Retryable() != item.retryable {
			test.Fatalf("%q retryable = %v", item.code, item.code.Retryable())
		}
	}
	if len(cases) != 9 {
		test.Fatalf("expected the nine stable codes, got %d", len(cases))
	}
}

func TestUnknown(test *testing.T) {
	test.Parallel()

	if code.Code("OTHER").Valid() {
		test.Fatal("an unknown code was accepted")
	}
}
