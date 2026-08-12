package wiring

import (
	"strings"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// brand carries what a hook writer needs: which agent is being registered, and the resolved Drizz command each hook
// runs.
type brand struct {
	agent   agent.Agent
	command string
}

// token is the recognizable middle of every Drizz hook command — "hook <agent>" — so a Drizz registration can be found
// and removed without disturbing a person's own hooks.
func (brand brand) token() string {
	return "hook " + brand.agent.Kind().String()
}

// invocation is the shell command a nested-style hook runs: the Drizz program, the hidden hook verb, the agent, and
// the slot, so the receiver knows who fired and what happened.
func (brand brand) invocation(slot agent.Slot) string {
	return brand.command + " " + brand.token() + " " + slot.String()
}

// editing is one hook write in progress: the document to modify and the registration to apply.
type editing struct {
	document map[string]any
	mark     brand
}

// stylist writes, removes, and detects Drizz's hook registration in a document. Each style is a different layout, so
// the layout is chosen once by style rather than branched on everywhere.
type stylist interface {
	inscribe(edit editing) error
	erase(document map[string]any)
	present(document map[string]any) bool
}

func (Store) stylist(style agent.Style) (stylist, error) {
	switch style {
	case agent.Claude:
		return nested{events: map[agent.Slot]string{agent.Prompt: "UserPromptSubmit", agent.Turn: "Stop"}}, nil
	case agent.Gemini:
		return nested{events: map[agent.Slot]string{agent.Prompt: "BeforeAgent", agent.Turn: "AfterAgent"}}, nil
	case agent.Codex:
		return notify{}, nil
	case agent.None:
		return nil, connect.Unsupported{}
	default:
		return nil, connect.Unsupported{}
	}
}

// nested writes the Claude and Gemini layout: a hooks object keyed by event, each event holding an array of handler
// groups. Drizz's group is appended per event, and only Drizz's group is removed, so a person's own hooks are kept.
type nested struct {
	events map[agent.Slot]string
}

func (nested nested) inscribe(edit editing) error {
	hooks, _ := edit.document["hooks"].(map[string]any)
	if hooks == nil {
		hooks = map[string]any{}
	}
	for slot, event := range nested.events {
		groups := nested.strip(hooks[event])
		groups = append(groups, map[string]any{
			"hooks": []any{
				map[string]any{"type": "command", "command": edit.mark.invocation(slot)},
			},
		})
		hooks[event] = groups
	}
	edit.document["hooks"] = hooks
	return nil
}

func (nested nested) erase(document map[string]any) {
	hooks, _ := document["hooks"].(map[string]any)
	if hooks == nil {
		return
	}
	for _, event := range nested.events {
		groups := nested.strip(hooks[event])
		if len(groups) == 0 {
			delete(hooks, event)
		} else {
			hooks[event] = groups
		}
	}
	if len(hooks) == 0 {
		delete(document, "hooks")
	} else {
		document["hooks"] = hooks
	}
}

func (nested nested) present(document map[string]any) bool {
	hooks, _ := document["hooks"].(map[string]any)
	for _, event := range nested.events {
		if array, valid := hooks[event].([]any); valid && nested.marked(array) {
			return true
		}
	}
	return false
}

// strip returns the handler groups for one event with any Drizz group removed.
func (nested nested) strip(value any) []any {
	array, _ := value.([]any)
	kept := make([]any, 0, len(array))
	for _, item := range array {
		if !nested.group(item) {
			kept = append(kept, item)
		}
	}
	return kept
}

// group reports whether a handler group is a Drizz registration, by looking for the hook token in any of its commands.
func (nested) group(value any) bool {
	entry, _ := value.(map[string]any)
	handlers, _ := entry["hooks"].([]any)
	for _, handler := range handlers {
		shape, _ := handler.(map[string]any)
		command, _ := shape["command"].(string)
		if strings.Contains(command, "hook ") {
			return true
		}
	}
	return false
}

// marked reports whether any handler group in the array is a Drizz registration.
func (nested nested) marked(array []any) bool {
	for _, item := range array {
		if nested.group(item) {
			return true
		}
	}
	return false
}

// notify writes the Codex layout: a single top-level notify program invoked when a turn completes. Codex allows only
// one such program, so Drizz refuses to overwrite a different one already there.
type notify struct{}

func (notify notify) inscribe(edit editing) error {
	if existing, present := edit.document["notify"]; present && !notify.command(existing) {
		return connect.Occupied{}
	}
	edit.document["notify"] = []any{edit.mark.command, "hook", edit.mark.agent.Kind().String(), agent.Turn.String()}
	return nil
}

func (notify notify) erase(document map[string]any) {
	if existing, present := document["notify"]; present && notify.command(existing) {
		delete(document, "notify")
	}
}

func (notify notify) present(document map[string]any) bool {
	return notify.command(document["notify"])
}

// command reports whether a notify value is a Drizz program, by looking for the hook verb among its arguments.
func (notify) command(value any) bool {
	array, _ := value.([]any)
	for _, item := range array {
		if word, _ := item.(string); word == "hook" {
			return true
		}
	}
	return false
}
