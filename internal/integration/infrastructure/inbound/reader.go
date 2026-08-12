// Package inbound reads an agent's hook event and turns it into a Drizz intake event. It knows the two ways agents
// deliver an event — as JSON on standard input, or as a single JSON argument — and the field names different agents
// use for a person's prompt and the agent's own reply, keeping that vendor detail out of the application core.
package inbound

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/DrizzDev/platform/internal/integration/application/intake"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// ceiling bounds how much of a hook event Drizz reads, so a runaway or hostile event cannot exhaust memory.
const ceiling = 1 << 20

// Reader reads one hook event from the channel the calling agent uses.
type Reader struct {
	source io.Reader
}

func New(source io.Reader) Reader {
	return Reader{source: source}
}

// Request is one event to read: which agent fired it, which turn moment it marks, and the argument payload present
// when the agent delivers the event as a single argument rather than on standard input.
type Request struct {
	Agent   agent.Agent
	Slot    agent.Slot
	Payload string
}

// Read produces the intake event for a hook notification. It never fails on content: an event it cannot parse yields
// an empty-text event, still recorded so the ordering of prompts, turns, and actions is preserved.
func (reader Reader) Read(request Request) intake.Event {
	return intake.Event{
		Agent: request.Agent.Kind(),
		Slot:  request.Slot,
		Text:  reader.extract(request),
	}
}

func (reader Reader) extract(request Request) string {
	document := reader.parse(request)
	for _, field := range reader.fields(request.Slot) {
		if text, valid := document[field].(string); valid && text != "" {
			return text
		}
	}
	if request.Slot == agent.Prompt {
		return reader.join(document["input-messages"])
	}
	return ""
}

func (reader Reader) parse(request Request) map[string]any {
	raw := []byte(request.Payload)
	if request.Agent.Hooking().Channel() == agent.Stdin {
		read, failure := io.ReadAll(io.LimitReader(reader.source, ceiling))
		if failure != nil {
			return map[string]any{}
		}
		raw = read
	}
	document := map[string]any{}
	if json.Unmarshal(raw, &document) != nil {
		return map[string]any{}
	}
	return document
}

// fields are the candidate keys, most specific first, that carry the human text for a turn moment across the agents
// Drizz supports.
func (Reader) fields(slot agent.Slot) []string {
	if slot == agent.Turn {
		return []string{"last-assistant-message", "response"}
	}
	return []string{"prompt", "last-user-message"}
}

// join renders an array of message strings as one block, so an agent that exposes the turn's prompts as a list still
// yields readable text.
func (Reader) join(value any) string {
	array, valid := value.([]any)
	if !valid {
		return ""
	}
	lines := make([]string, 0, len(array))
	for _, item := range array {
		if text, valid := item.(string); valid {
			lines = append(lines, text)
		}
	}
	return strings.Join(lines, "\n")
}
