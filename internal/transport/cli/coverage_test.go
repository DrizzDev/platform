package cli_test

import (
	"io"
	"testing"

	"github.com/DrizzDev/platform/internal/capability/domain/catalog"
	"github.com/DrizzDev/platform/internal/transport/cli"
)

// TestCatalogCovered is half the completeness gate: every catalogued capability must have a command-line command. With
// no central dispatch, this is what guarantees a capability is never added to the catalogue but left off this surface.
func TestCatalogCovered(test *testing.T) {
	test.Parallel()

	names, failure := cli.Names(fixture{output: io.Discard}.options())
	if failure != nil {
		test.Fatal(failure)
	}
	present := make(map[string]bool, len(names))
	for _, name := range names {
		present[name] = true
	}
	for _, entry := range catalog.New().List() {
		if !present[cli.Slug(entry.Title())] {
			test.Errorf("catalogued capability %q has no command-line command", entry.Title())
		}
	}
}
