package sqlite_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DrizzDev/platform/internal/identity/infrastructure/sqlite"
)

func TestLocate(test *testing.T) {
	test.Parallel()

	path, failure := (sqlite.Location{Root: test.TempDir()}).Resolve()
	if failure != nil {
		test.Fatal(failure)
	}
	if filepath.Base(path) != "identity.db" || filepath.Base(filepath.Dir(path)) != "drizz" {
		test.Fatalf("path = %q", path)
	}
	info, failure := os.Stat(filepath.Dir(path))
	if failure != nil {
		test.Fatal(failure)
	}
	if info.Mode().Perm() != 0o700 {
		test.Fatalf("directory permissions = %o", info.Mode().Perm())
	}

	fixture{test: test}.open(path)
}

func TestSymlink(test *testing.T) {
	test.Parallel()

	root := test.TempDir()
	if failure := os.Symlink(test.TempDir(), filepath.Join(root, "drizz")); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := (sqlite.Location{Root: root}).Resolve(); failure == nil {
		test.Fatal("a symlinked data directory was accepted")
	}
}

func TestFileSymlink(test *testing.T) {
	test.Parallel()

	root := test.TempDir()
	directory := filepath.Join(root, "drizz")
	if failure := os.MkdirAll(directory, 0o700); failure != nil {
		test.Fatal(failure)
	}
	if failure := os.Symlink(filepath.Join(test.TempDir(), "elsewhere.db"), filepath.Join(directory, "identity.db")); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := (sqlite.Location{Root: root}).Resolve(); failure == nil {
		test.Fatal("a symlinked database file was accepted")
	}
}
