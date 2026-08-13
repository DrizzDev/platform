package wiring_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandInstallsAndRemoves(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)

	conflict, failure := kit.store().Command(context.Background(), kit.job(kit.pick("claude-code")))
	if failure != nil {
		test.Fatal(failure)
	}
	if conflict {
		test.Fatal("a fresh install must not report a conflict")
	}

	path := filepath.Join(home, ".claude", "commands", "author.md")
	raw, failure := os.ReadFile(path)
	if failure != nil {
		test.Fatalf("command file was not written: %v", failure)
	}
	for _, want := range []string{"drizz:managed", "$ARGUMENTS", "Name the role"} {
		if !strings.Contains(string(raw), want) {
			test.Fatalf("command file is missing %q", want)
		}
	}

	if failure := kit.store().Uncommand(context.Background(), kit.pick("claude-code")); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := os.Stat(path); !os.IsNotExist(failure) {
		test.Fatal("Uncommand left the Drizz command file behind")
	}
}

func TestCommandKeepsForeignFile(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)

	path := filepath.Join(home, ".claude", "commands", "author.md")
	if failure := os.MkdirAll(filepath.Dir(path), 0o700); failure != nil {
		test.Fatal(failure)
	}
	mine := "my own author command\n"
	if failure := os.WriteFile(path, []byte(mine), 0o600); failure != nil {
		test.Fatal(failure)
	}

	conflict, failure := kit.store().Command(context.Background(), kit.job(kit.pick("claude-code")))
	if failure != nil {
		test.Fatal(failure)
	}
	if !conflict {
		test.Fatal("a command file Drizz does not own must be reported as a conflict")
	}
	if raw, _ := os.ReadFile(path); string(raw) != mine {
		test.Fatalf("Command overwrote a foreign file: %q", raw)
	}

	if failure := kit.store().Uncommand(context.Background(), kit.pick("claude-code")); failure != nil {
		test.Fatal(failure)
	}
	if _, failure := os.Stat(path); failure != nil {
		test.Fatal("Uncommand removed a command file Drizz does not own")
	}
}
