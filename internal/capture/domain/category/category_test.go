package category_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capture/domain/category"
)

func TestCategory(test *testing.T) {
	test.Parallel()

	known := []category.Category{
		category.Prompt, category.Response, category.Tool,
		category.Screen, category.Hierarchy, category.File, category.Log,
	}
	for _, value := range known {
		if !value.Valid() {
			test.Fatalf("category %q rejected", value)
		}
	}
	if category.Category("SECRET").Valid() {
		test.Fatal("an unknown category was accepted")
	}
	if category.Prompt.String() != "PROMPT" {
		test.Fatalf("category string = %q", category.Prompt.String())
	}
}
