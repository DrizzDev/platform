package sqlite_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/DrizzDev/platform/internal/capture/infrastructure/sqlite"
)

func TestMigrate(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.open(filepath.Join(test.TempDir(), "capture.db"))
	checks := map[string]string{
		"PRAGMA user_version": "1",
		"PRAGMA journal_mode": "wal",
		"PRAGMA foreign_keys": "1",
		"SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'journal'": "1",
	}
	for statement, expected := range checks {
		value, failure := store.Query(statement)
		if failure != nil {
			test.Fatalf("%s: %v", statement, failure)
		}
		if value != expected {
			test.Fatalf("%s = %q, want %q", statement, value, expected)
		}
	}
}

func TestIdempotent(test *testing.T) {
	test.Parallel()

	kit := fixture{test: test}
	path := filepath.Join(test.TempDir(), "capture.db")
	kit.open(path)

	reopened := kit.open(path)
	value, failure := reopened.Query("PRAGMA user_version")
	if failure != nil {
		test.Fatal(failure)
	}
	if value != "1" {
		test.Fatalf("reopened version = %q", value)
	}
}

func TestCorrupt(test *testing.T) {
	test.Parallel()

	path := filepath.Join(test.TempDir(), "capture.db")
	if failure := os.WriteFile(path, []byte("this is not a database"), 0o600); failure != nil {
		test.Fatal(failure)
	}
	store := fixture{test: test}.open(path)
	value, failure := store.Query("PRAGMA user_version")
	if failure != nil {
		test.Fatalf("recovered store is unusable: %v", failure)
	}
	if value != "1" {
		test.Fatalf("recovered version = %q", value)
	}
	aside, _ := filepath.Glob(path + ".corrupt-*")
	if len(aside) == 0 {
		test.Fatal("the corrupt database was not quarantined")
	}
}

func TestDependencies(test *testing.T) {
	test.Parallel()

	if _, failure := sqlite.New(context.Background(), sqlite.Options{}); failure == nil {
		test.Fatal("a database without options was accepted")
	}
}

func TestStorage(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.open(filepath.Join(test.TempDir(), "capture.db"))
	pages, failure := store.Query("PRAGMA page_count")
	if failure != nil {
		test.Fatal(failure)
	}
	if failure := store.Exec("PRAGMA max_page_count = " + pages); failure != nil {
		test.Fatal(failure)
	}
	insert := "INSERT INTO journal (trace, span, parent, sequence, mark, origin, fidelity, category, payload, state, stamped) " +
		"VALUES ('t', 's', '', 0, 'EXACT', 'CAPABILITY', 'EXACT', 'TOOL', randomblob(300000), 'PENDING', 0)"
	if failure := store.Exec(insert); failure == nil {
		test.Fatal("a write past the page ceiling was accepted")
	}
}
