package wiring

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// hookpath resolves the agent's hook configuration file from its hook descriptor. It may be the same file the agent
// keeps its MCP servers in, or a different one.
func (store Store) hookpath(hooking agent.Hooking) string {
	root, failure := store.anchor(hooking.Base())
	if failure != nil {
		return ""
	}
	return filepath.Join(append([]string{root}, hooking.Segments()...)...)
}

// load reads the agent's hook configuration into a generic document, in the descriptor's dialect. An absent file is an
// empty document; an unreadable file is Locked; a file that will not parse is Malformed.
func (store Store) load(hooking agent.Hooking) (map[string]any, error) {
	path := store.hookpath(hooking)
	if path == "" {
		return nil, connect.Locked{}
	}
	if failure := store.screen(path); failure != nil {
		return nil, failure
	}
	raw, failure := os.ReadFile(path)
	if errors.Is(failure, os.ErrNotExist) {
		return map[string]any{}, nil
	}
	if failure != nil {
		return nil, connect.Locked{}
	}
	coder, failure := store.codec(hooking.Dialect())
	if failure != nil {
		return nil, failure
	}
	document, failure := coder.parse(raw)
	if failure != nil {
		return nil, connect.Malformed{}
	}
	return document, nil
}
