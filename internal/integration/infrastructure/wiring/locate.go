package wiring

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// locate resolves the agent's configuration file path from its descriptor and refuses a path whose file is a symlink
// or is otherwise not a regular file, so a planted node cannot redirect a write elsewhere.
func (store Store) locate(target agent.Agent) (string, error) {
	root, failure := store.anchor(target.Base())
	if failure != nil {
		return "", failure
	}
	path := filepath.Join(append([]string{root}, target.Segments()...)...)
	if failure := store.screen(path); failure != nil {
		return "", failure
	}
	return path, nil
}

// anchor turns a descriptor's base into a real directory for the running user without the domain ever making an
// operating-system call.
func (Store) anchor(base agent.Base) (string, error) {
	switch base {
	case agent.Home:
		return os.UserHomeDir()
	case agent.Config:
		return os.UserConfigDir()
	default:
		return "", connect.Locked{}
	}
}

func (Store) screen(path string) error {
	info, failure := os.Lstat(path)
	switch {
	case errors.Is(failure, os.ErrNotExist):
		return nil
	case failure != nil:
		return connect.Locked{}
	case info.Mode()&os.ModeSymlink != 0:
		return connect.Locked{}
	case !info.Mode().IsRegular():
		return connect.Locked{}
	}
	return nil
}

func (Store) present(path string) bool {
	_, failure := os.Stat(path)
	return failure == nil
}

// read loads the agent's configuration into a generic document. An absent file is an empty document, not an error; an
// unreadable file is Locked; a file that will not parse is Malformed, so the installer never overwrites a file it
// cannot understand.
func (store Store) read(target agent.Agent) (map[string]any, error) {
	path, failure := store.locate(target)
	if failure != nil {
		return nil, failure
	}
	raw, failure := os.ReadFile(path)
	if errors.Is(failure, os.ErrNotExist) {
		return map[string]any{}, nil
	}
	if failure != nil {
		return nil, connect.Locked{}
	}
	coder, failure := store.codec(target.Dialect())
	if failure != nil {
		return nil, failure
	}
	document, failure := coder.parse(raw)
	if failure != nil {
		return nil, connect.Malformed{}
	}
	return document, nil
}
