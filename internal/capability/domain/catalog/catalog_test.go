package catalog_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
)

func TestList(test *testing.T) {
	test.Parallel()

	entries := catalog.New().List()
	if len(entries) == 0 {
		test.Fatal("the catalog is empty")
	}
	seen := make(map[string]bool, len(entries))
	for _, entry := range entries {
		if entry.Name() == "" || entry.Title() == "" || entry.Summary() == "" {
			test.Errorf("capability %q is missing a name, title, or summary", entry.Name())
		}
		if seen[entry.Name()] {
			test.Errorf("capability %q is listed more than once", entry.Name())
		}
		seen[entry.Name()] = true
	}
}

func TestLookup(test *testing.T) {
	test.Parallel()

	shot, found := catalog.New().Lookup(catalog.Screenshot)
	if !found {
		test.Fatal("screenshot capability is missing from the catalog")
	}
	if shot.Summary() == "" {
		test.Fatal("screenshot capability has no summary")
	}
	parameters := shot.Parameters()
	if len(parameters) != 1 || parameters[0].Name() != "serial" {
		test.Fatalf("screenshot parameters = %v, want one named serial", parameters)
	}
}

func TestUnknown(test *testing.T) {
	test.Parallel()

	if _, found := catalog.New().Lookup("nonexistent"); found {
		test.Fatal("an unknown capability must not be found")
	}
}
