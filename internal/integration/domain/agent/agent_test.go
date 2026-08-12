package agent_test

import (
	"testing"

	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

func TestCatalogLooksUpAgents(test *testing.T) {
	test.Parallel()

	catalog := agent.New()
	if len(catalog.List()) == 0 {
		test.Fatal("the agent catalog must not be empty")
	}

	codex, found := catalog.Lookup("codex")
	if !found {
		test.Fatal("codex is missing from the catalog")
	}
	if codex.Dialect() != agent.Toml {
		test.Fatalf("codex dialect = %s, want TOML", codex.Dialect())
	}
	if codex.Collection() != "mcp_servers" {
		test.Fatalf("codex collection = %s, want mcp_servers", codex.Collection())
	}

	if _, found := catalog.Lookup("nonesuch"); found {
		test.Fatal("an unknown agent was accepted")
	}
}

func TestClaudeCodeCarriesStdioType(test *testing.T) {
	test.Parallel()

	target, found := agent.New().Lookup("claude-code")
	if !found {
		test.Fatal("claude-code is missing from the catalog")
	}
	if !target.Typed() {
		test.Fatal("claude-code must declare an explicit stdio type")
	}
	if target.Dialect() != agent.Json {
		test.Fatalf("claude-code dialect = %s, want JSON", target.Dialect())
	}
}
