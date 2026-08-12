package wiring_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaptureWritesClaudeHooks(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)

	if failure := kit.store().Capture(context.Background(), kit.job(kit.pick("claude-code"))); failure != nil {
		test.Fatal(failure)
	}

	document := kit.read(filepath.Join(home, ".claude", "settings.json"))
	hooks, valid := document["hooks"].(map[string]any)
	if !valid {
		test.Fatal("the hooks object was not written")
	}
	for _, event := range []string{"UserPromptSubmit", "Stop"} {
		groups, valid := hooks[event].([]any)
		if !valid || len(groups) == 0 {
			test.Fatalf("no hook registered for %s", event)
		}
	}
	// The MCP configuration file must be left alone — hooks live in a different file.
	if _, failure := os.Stat(filepath.Join(home, ".claude.json")); failure == nil {
		test.Fatal("capture must not touch the MCP configuration file")
	}
}

func TestCapturePreservesOtherHooks(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".claude", "settings.json")
	kit.seed(file{
		path:    path,
		content: `{"hooks":{"UserPromptSubmit":[{"hooks":[{"type":"command","command":"mine.sh"}]}]}}`,
	})

	if failure := kit.store().Capture(context.Background(), kit.job(kit.pick("claude-code"))); failure != nil {
		test.Fatal(failure)
	}

	groups := kit.read(path)["hooks"].(map[string]any)["UserPromptSubmit"].([]any)
	if len(groups) != 2 {
		test.Fatalf("expected the person's hook plus Drizz's, got %d groups", len(groups))
	}
	if !strings.Contains(kit.text(path), "mine.sh") {
		test.Fatal("the person's own hook was lost")
	}
}

func TestUncaptureRemovesOnlyDrizzHooks(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".claude", "settings.json")
	kit.seed(file{
		path:    path,
		content: `{"hooks":{"UserPromptSubmit":[{"hooks":[{"type":"command","command":"mine.sh"}]}]}}`,
	})

	store := kit.store()
	if failure := store.Capture(context.Background(), kit.job(kit.pick("claude-code"))); failure != nil {
		test.Fatal(failure)
	}
	if failure := store.Uncapture(context.Background(), kit.pick("claude-code")); failure != nil {
		test.Fatal(failure)
	}

	text := kit.text(path)
	if strings.Contains(text, "hook claude-code") {
		test.Fatal("Drizz's hook was not removed")
	}
	if !strings.Contains(text, "mine.sh") {
		test.Fatal("the person's own hook was removed")
	}
}

func TestCaptureCodexNotify(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".codex", "config.toml")
	kit.seed(file{path: path, content: "model = \"o3\"\n"})

	store := kit.store()
	if failure := store.Capture(context.Background(), kit.job(kit.pick("codex"))); failure != nil {
		test.Fatal(failure)
	}

	captures, failure := store.Captures(kit.pick("codex"))
	if failure != nil || !captures {
		test.Fatalf("codex notify was not registered: captures=%v err=%v", captures, failure)
	}
	document := kit.decode(path)
	if document["model"] != "o3" {
		test.Fatal("an unrelated Codex setting was lost")
	}
	if _, present := document["notify"]; !present {
		test.Fatal("the Codex notify program was not written")
	}
}

func TestCaptureCodexRefusesExistingNotify(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".codex", "config.toml")
	kit.seed(file{path: path, content: "notify = [\"my-notifier\"]\n"})

	failure := kit.store().Capture(context.Background(), kit.job(kit.pick("codex")))
	if failure == nil {
		test.Fatal("capture must refuse to overwrite a different Codex notify program")
	}
	if !strings.Contains(kit.text(path), "my-notifier") {
		test.Fatal("the person's own notify program was overwritten")
	}
}

func TestCaptureUnsupportedAgent(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)

	failure := kit.store().Capture(context.Background(), kit.job(kit.pick("claude-desktop")))
	if failure == nil {
		test.Fatal("capture must report that Claude Desktop has no hook mechanism")
	}
}

func (kit kit) text(path string) string {
	kit.test.Helper()
	return string(kit.slurp(path))
}
