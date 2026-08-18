package binary_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/integration/infrastructure/binary"
)

// TestLocateKeepsLauncherSymlink reproduces the Homebrew layout: a stable launcher symlink (bin/drizz) pointing at a
// version-pinned target (Caskroom/<version>/drizz). The entry the installer writes must record the launcher, not the
// target, so it survives an upgrade that replaces the target.
func TestLocateKeepsLauncherSymlink(test *testing.T) {
	test.Parallel()

	root := test.TempDir()

	cask := filepath.Join(root, "Caskroom", "drizz", "0.1.5")
	if failure := os.MkdirAll(cask, 0o755); failure != nil {
		test.Fatal(failure)
	}
	target := filepath.Join(cask, "drizz")
	if failure := os.WriteFile(target, []byte("binary"), 0o755); failure != nil {
		test.Fatal(failure)
	}

	bin := filepath.Join(root, "bin")
	if failure := os.MkdirAll(bin, 0o755); failure != nil {
		test.Fatal(failure)
	}
	launcher := filepath.Join(bin, "drizz")
	if failure := os.Symlink(target, launcher); failure != nil {
		test.Fatal(failure)
	}

	resolver := binary.Using(func() (string, error) { return launcher, nil })
	located, failure := resolver.Locate()
	if failure != nil {
		test.Fatal(failure)
	}
	if located != launcher {
		test.Fatalf("Locate followed the launcher to a version-pinned target: got %q, want %q", located, launcher)
	}
	if strings.Contains(located, "Caskroom") {
		test.Fatalf("Locate returned a version-pinned path that breaks on upgrade: %q", located)
	}
}
