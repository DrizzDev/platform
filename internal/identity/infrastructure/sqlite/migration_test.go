package sqlite_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInterrupted(test *testing.T) {
	test.Parallel()

	path := filepath.Join(test.TempDir(), "identity.db")
	fixture := fixture{test: test}
	store := fixture.open(path)

	if failure := store.Exec("PRAGMA user_version = 0"); failure != nil {
		test.Fatal(failure)
	}
	if failure := store.Close(); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := fixture.build(path); failure == nil {
		test.Fatal("an inconsistent prior migration was not rejected")
	}
}

func TestConstraints(test *testing.T) {
	test.Parallel()

	store := fixture{test: test}.open(filepath.Join(test.TempDir(), "identity.db"))
	cases := map[string]string{
		"reason":   "INSERT INTO cleanup (name, reason, state, attempts, next, deadline, created) VALUES ('k', 'BOGUS', 'PENDING', 0, 0, 0, 0)",
		"state":    "INSERT INTO cleanup (name, reason, state, attempts, next, deadline, created) VALUES ('k', 'LOGOUT', 'BOGUS', 0, 0, 0, 0)",
		"revision": "INSERT INTO pointer (session, name, revision, epoch) VALUES ('s', 'k', 0, 0)",
		"expiry":   "INSERT INTO lease (session, owner, expiry) VALUES ('s', 'o', 0)",
	}
	for name, statement := range cases {
		test.Run(name, func(test *testing.T) {
			test.Parallel()
			if failure := store.Exec(statement); failure == nil {
				test.Fatal("the database accepted a value the constraint forbids")
			}
		})
	}
}

func FuzzOpen(fuzz *testing.F) {
	fuzz.Add([]byte("this is not a database"))
	fuzz.Add([]byte{0})
	fuzz.Fuzz(func(test *testing.T, data []byte) {
		path := filepath.Join(test.TempDir(), "identity.db")
		if failure := os.WriteFile(path, data, 0o600); failure != nil {
			test.Fatal(failure)
		}
		store, failure := fixture{test: test}.build(path)
		if failure == nil {
			_ = store.Close()
		}
	})
}
