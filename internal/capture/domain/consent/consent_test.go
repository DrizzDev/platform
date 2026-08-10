package consent_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/category"
	"github.com/DrizzDev/platform/internal/capture/domain/consent"
)

func TestConsent(test *testing.T) {
	test.Parallel()

	value, failure := consent.New(consent.Input{Approved: []category.Category{category.Prompt, category.Screen}})
	if failure != nil {
		test.Fatal(failure)
	}
	if !value.Allows(category.Prompt) || !value.Allows(category.Screen) {
		test.Fatal("an approved category was not allowed")
	}
	if value.Allows(category.Log) {
		test.Fatal("an unapproved category was allowed")
	}
}

func TestConsentRejectsUnknown(test *testing.T) {
	test.Parallel()

	if _, failure := consent.New(consent.Input{Approved: []category.Category{category.Category("SECRET")}}); failure == nil {
		test.Fatal("consent for an unknown category was accepted")
	}
}
