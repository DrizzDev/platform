package wiring_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
	"github.com/DrizzDev/platform/internal/integration/domain/server"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/wiring"
)

type kit struct {
	test *testing.T
}

// store builds a wiring store over no-op telemetry, so a test exercises real file merging without a telemetry backend.
func (kit kit) store() wiring.Store {
	kit.test.Helper()
	made, failure := wiring.New(wiring.Options{
		Tracer: tracenoop.NewTracerProvider().Tracer("test"),
		Meter:  metricnoop.NewMeterProvider().Meter("test"),
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return made
}

// file is one seeded configuration file: where it lives and what it contains.
type file struct {
	path    string
	content string
}

func (kit kit) pick(kind string) agent.Agent {
	kit.test.Helper()
	target, found := agent.New().Lookup(agent.Kind(kind))
	if !found {
		kit.test.Fatalf("agent %q is missing from the catalog", kind)
	}
	return target
}

func (kit kit) job(target agent.Agent) connect.Task {
	kit.test.Helper()
	entry, failure := server.New(server.Input{
		Name:    connect.Name,
		Command: "/opt/drizz/drizz",
		Args:    []string{connect.Launch},
	})
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return connect.Task{Agent: target, Server: entry}
}

func (kit kit) read(path string) map[string]any {
	kit.test.Helper()
	var document map[string]any
	if failure := json.Unmarshal(kit.slurp(path), &document); failure != nil {
		kit.test.Fatalf("configuration is not valid JSON: %v", failure)
	}
	return document
}

func (kit kit) servers(path string) map[string]any {
	kit.test.Helper()
	group, valid := kit.read(path)["mcpServers"].(map[string]any)
	if !valid {
		kit.test.Fatal("configuration has no mcpServers object")
	}
	return group
}

func (kit kit) decode(path string) map[string]any {
	kit.test.Helper()
	var document map[string]any
	if failure := toml.Unmarshal(kit.slurp(path), &document); failure != nil {
		kit.test.Fatalf("configuration is not valid TOML: %v", failure)
	}
	return document
}

func (kit kit) slurp(path string) []byte {
	kit.test.Helper()
	raw, failure := os.ReadFile(path)
	if failure != nil {
		kit.test.Fatal(failure)
	}
	return raw
}

func (kit kit) seed(item file) {
	kit.test.Helper()
	if failure := os.MkdirAll(filepath.Dir(item.path), 0o700); failure != nil {
		kit.test.Fatal(failure)
	}
	if failure := os.WriteFile(item.path, []byte(item.content), 0o600); failure != nil {
		kit.test.Fatal(failure)
	}
}

func TestConnectCreatesEntry(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)

	if failure := kit.store().Connect(context.Background(), kit.job(kit.pick("claude-code"))); failure != nil {
		test.Fatal(failure)
	}

	drizz, valid := kit.servers(filepath.Join(home, ".claude.json"))["drizz"].(map[string]any)
	if !valid {
		test.Fatal("the Drizz entry was not written")
	}
	if drizz["command"] != "/opt/drizz/drizz" {
		test.Fatalf("command = %v", drizz["command"])
	}
	if drizz["type"] != "stdio" {
		test.Fatalf("claude-code entry must carry a stdio type, got %v", drizz["type"])
	}
}

func TestConnectOmitsTypeWhenNotRequired(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)

	if failure := kit.store().Connect(context.Background(), kit.job(kit.pick("gemini"))); failure != nil {
		test.Fatal(failure)
	}

	drizz := kit.servers(filepath.Join(home, ".gemini", "settings.json"))["drizz"].(map[string]any)
	if _, present := drizz["type"]; present {
		test.Fatal("gemini entry must not carry a type field")
	}
}

func TestConnectPreservesOtherSettings(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".claude.json")
	kit.seed(file{path: path, content: `{"theme":"dark","mcpServers":{"weather":{"command":"weather"}}}`})

	if failure := kit.store().Connect(context.Background(), kit.job(kit.pick("claude-code"))); failure != nil {
		test.Fatal(failure)
	}

	document := kit.read(path)
	if document["theme"] != "dark" {
		test.Fatalf("an unrelated setting was lost: theme = %v", document["theme"])
	}
	group := document["mcpServers"].(map[string]any)
	if _, present := group["weather"]; !present {
		test.Fatal("an unrelated MCP server was lost")
	}
	if _, present := group["drizz"]; !present {
		test.Fatal("the Drizz entry was not added")
	}
}

func TestConnectIsIdempotent(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	store := kit.store()
	target := kit.pick("claude-code")

	for range 2 {
		if failure := store.Connect(context.Background(), kit.job(target)); failure != nil {
			test.Fatal(failure)
		}
	}

	group := kit.servers(filepath.Join(home, ".claude.json"))
	if len(group) != 1 {
		test.Fatalf("expected exactly one server entry, got %d", len(group))
	}
}

func TestConnectRefusesMalformed(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".claude.json")
	kit.seed(file{path: path, content: `{ this is not json `})

	failure := kit.store().Connect(context.Background(), kit.job(kit.pick("claude-code")))
	if _, malformed := errors.AsType[connect.Malformed](failure); !malformed {
		test.Fatalf("a malformed configuration must be refused, got %v", failure)
	}
	if string(kit.slurp(path)) != `{ this is not json ` {
		test.Fatal("a refused configuration must be left untouched")
	}
}

func TestConnectBacksUpOriginal(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".claude.json")
	kit.seed(file{path: path, content: `{"theme":"dark"}`})

	if failure := kit.store().Connect(context.Background(), kit.job(kit.pick("claude-code"))); failure != nil {
		test.Fatal(failure)
	}
	if string(kit.slurp(path+".drizz.bak")) != `{"theme":"dark"}` {
		test.Fatal("the original configuration was not backed up before the change")
	}
}

func TestDisconnectRemovesOnlyDrizz(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".claude.json")
	kit.seed(file{path: path, content: `{"mcpServers":{"weather":{"command":"weather"},"drizz":{"command":"old"}}}`})

	if failure := kit.store().Disconnect(context.Background(), kit.pick("claude-code")); failure != nil {
		test.Fatal(failure)
	}

	group := kit.servers(path)
	if _, present := group["drizz"]; present {
		test.Fatal("the Drizz entry was not removed")
	}
	if _, present := group["weather"]; !present {
		test.Fatal("an unrelated MCP server was removed")
	}
}

func TestDisconnectWhenAbsentSucceeds(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	kit.seed(file{path: filepath.Join(home, ".claude.json"), content: `{"theme":"dark"}`})

	if failure := kit.store().Disconnect(context.Background(), kit.pick("claude-code")); failure != nil {
		test.Fatalf("removing an absent entry must succeed, got %v", failure)
	}
}

func TestConnectRejectsSymlink(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".claude.json")
	kit.seed(file{path: filepath.Join(home, "real.json"), content: `{}`})
	if failure := os.Symlink(filepath.Join(home, "real.json"), path); failure != nil {
		test.Fatal(failure)
	}

	failure := kit.store().Connect(context.Background(), kit.job(kit.pick("claude-code")))
	if _, locked := errors.AsType[connect.Locked](failure); !locked {
		test.Fatalf("a symlinked configuration path must be refused, got %v", failure)
	}
}

func TestConnectCodexToml(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".codex", "config.toml")
	kit.seed(file{path: path, content: "model = \"o3\"\n\n[mcp_servers.weather]\ncommand = \"weather\"\n"})

	if failure := kit.store().Connect(context.Background(), kit.job(kit.pick("codex"))); failure != nil {
		test.Fatal(failure)
	}

	document := kit.decode(path)
	if document["model"] != "o3" {
		test.Fatalf("an unrelated Codex setting was lost: model = %v", document["model"])
	}
	group, valid := document["mcp_servers"].(map[string]any)
	if !valid {
		test.Fatal("the Codex mcp_servers table is missing")
	}
	if _, present := group["weather"]; !present {
		test.Fatal("an unrelated Codex MCP server was lost")
	}
	drizz, valid := group["drizz"].(map[string]any)
	if !valid {
		test.Fatal("the Drizz entry was not written into the Codex configuration")
	}
	if drizz["command"] != "/opt/drizz/drizz" {
		test.Fatalf("command = %v", drizz["command"])
	}
}

func TestDisconnectCodexToml(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	path := filepath.Join(home, ".codex", "config.toml")
	kit.seed(file{path: path, content: "[mcp_servers.weather]\ncommand = \"weather\"\n\n[mcp_servers.drizz]\ncommand = \"old\"\n"})

	if failure := kit.store().Disconnect(context.Background(), kit.pick("codex")); failure != nil {
		test.Fatal(failure)
	}

	group := kit.decode(path)["mcp_servers"].(map[string]any)
	if _, present := group["drizz"]; present {
		test.Fatal("the Drizz entry was not removed from the Codex configuration")
	}
	if _, present := group["weather"]; !present {
		test.Fatal("an unrelated Codex MCP server was removed")
	}
}

func TestDetectByDirectory(test *testing.T) {
	kit := kit{test: test}
	home := test.TempDir()
	test.Setenv("HOME", home)
	store := kit.store()
	target := kit.pick("codex")

	present, failure := store.Detect(target)
	if failure != nil || present {
		test.Fatalf("codex must be absent before its directory exists: present=%v err=%v", present, failure)
	}
	if failure := os.MkdirAll(filepath.Join(home, ".codex"), 0o700); failure != nil {
		test.Fatal(failure)
	}
	present, failure = store.Detect(target)
	if failure != nil || !present {
		test.Fatalf("codex must be detected once its directory exists: present=%v err=%v", present, failure)
	}
}
