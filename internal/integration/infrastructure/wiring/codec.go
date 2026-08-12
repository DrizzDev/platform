package wiring

import (
	"encoding/json"

	"github.com/pelletier/go-toml/v2"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// codec parses an agent configuration file into a generic document and renders a document back. The merge logic is
// identical across formats; only serialization differs by dialect, so a new format is one codec, not a new merge.
type codec interface {
	parse([]byte) (map[string]any, error)
	render(map[string]any) ([]byte, error)
}

func (Store) codec(dialect agent.Dialect) (codec, error) {
	switch dialect {
	case agent.Json:
		return jsoncodec{}, nil
	case agent.Toml:
		return tomlcodec{}, nil
	default:
		return nil, connect.Locked{}
	}
}

// entry is the generic shape of the Drizz server entry: the command and arguments that launch it, and — only for an
// agent that requires it — an explicit stdio type.
func (Store) entry(job connect.Task) map[string]any {
	shape := map[string]any{
		"command": job.Server.Command(),
		"args":    job.Server.Args(),
	}
	if job.Agent.Typed() {
		shape["type"] = "stdio"
	}
	return shape
}

type jsoncodec struct{}

func (jsoncodec) parse(raw []byte) (map[string]any, error) {
	document := map[string]any{}
	if len(raw) == 0 {
		return document, nil
	}
	if failure := json.Unmarshal(raw, &document); failure != nil {
		return nil, failure
	}
	return document, nil
}

func (jsoncodec) render(document map[string]any) ([]byte, error) {
	raw, failure := json.MarshalIndent(document, "", "  ")
	if failure != nil {
		return nil, failure
	}
	return append(raw, '\n'), nil
}

type tomlcodec struct{}

func (tomlcodec) parse(raw []byte) (map[string]any, error) {
	document := map[string]any{}
	if len(raw) == 0 {
		return document, nil
	}
	if failure := toml.Unmarshal(raw, &document); failure != nil {
		return nil, failure
	}
	return document, nil
}

func (tomlcodec) render(document map[string]any) ([]byte, error) {
	raw, failure := toml.Marshal(document)
	if failure != nil {
		return nil, failure
	}
	return raw, nil
}
