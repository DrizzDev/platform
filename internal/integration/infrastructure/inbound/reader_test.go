package inbound_test

import (
	"strings"
	"testing"

	"github.com/DrizzDev/platform/internal/integration/domain/agent"
	"github.com/DrizzDev/platform/internal/integration/infrastructure/inbound"
)

type kit struct {
	test *testing.T
}

func (kit kit) pick(kind string) agent.Agent {
	kit.test.Helper()
	target, found := agent.New().Lookup(agent.Kind(kind))
	if !found {
		kit.test.Fatalf("agent %q is missing from the catalog", kind)
	}
	return target
}

func TestReadStdinPrompt(test *testing.T) {
	test.Parallel()
	source := strings.NewReader(`{"prompt":"open the settings screen","cwd":"/x"}`)

	event := inbound.New(source).Read(inbound.Request{Agent: kit{test: test}.pick("claude-code"), Slot: agent.Prompt})

	if event.Text != "open the settings screen" {
		test.Fatalf("text = %q", event.Text)
	}
	if event.Agent != "claude-code" || event.Slot != agent.Prompt {
		test.Fatalf("event = %+v", event)
	}
}

func TestReadArgvTurn(test *testing.T) {
	test.Parallel()
	// Codex delivers the event as a single JSON argument, not on standard input.
	payload := `{"type":"agent-turn-complete","last-assistant-message":"done"}`

	event := inbound.New(strings.NewReader("")).Read(inbound.Request{
		Agent:   kit{test: test}.pick("codex"),
		Slot:    agent.Turn,
		Payload: payload,
	})

	if event.Text != "done" {
		test.Fatalf("text = %q", event.Text)
	}
}

func TestReadJoinsInputMessages(test *testing.T) {
	test.Parallel()
	payload := `{"input-messages":["first","second"]}`

	event := inbound.New(strings.NewReader("")).Read(inbound.Request{
		Agent:   kit{test: test}.pick("codex"),
		Slot:    agent.Prompt,
		Payload: payload,
	})

	if event.Text != "first\nsecond" {
		test.Fatalf("text = %q", event.Text)
	}
}

func TestReadUnparsableYieldsMarker(test *testing.T) {
	test.Parallel()
	event := inbound.New(strings.NewReader("not json at all")).Read(inbound.Request{
		Agent: kit{test: test}.pick("claude-code"),
		Slot:  agent.Prompt,
	})

	if event.Text != "" {
		test.Fatalf("an unparsable event must yield empty text, got %q", event.Text)
	}
	if event.Agent != "claude-code" {
		test.Fatal("an unparsable event must still be recorded as a marker for the agent")
	}
}
