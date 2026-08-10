package catalogue_test

import (
	"testing"
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/catalogue"
	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/policy"
)

type fixture struct {
	test *testing.T
}

func (fixture fixture) rule() policy.Policy {
	fixture.test.Helper()
	rule, failure := policy.New(policy.Input{Limit: 1024, Retention: time.Hour})
	if failure != nil {
		fixture.test.Fatal(failure)
	}
	return rule
}

func TestCatalogue(test *testing.T) {
	test.Parallel()

	book, failure := catalogue.New(catalogue.Input{Policies: map[category.Category]policy.Policy{
		category.Prompt: fixture{test: test}.rule(),
	}})
	if failure != nil {
		test.Fatal(failure)
	}
	if _, failure := book.Policy(category.Prompt); failure != nil {
		test.Fatalf("classified category rejected: %v", failure)
	}
	if _, failure := book.Policy(category.Screen); failure == nil {
		test.Fatal("an unclassified category did not fail closed")
	}
}

func TestCatalogueRejectsUnknown(test *testing.T) {
	test.Parallel()

	_, failure := catalogue.New(catalogue.Input{Policies: map[category.Category]policy.Policy{
		category.Category("SECRET"): fixture{test: test}.rule(),
	}})
	if failure == nil {
		test.Fatal("a catalogue with an unknown category was accepted")
	}
}
