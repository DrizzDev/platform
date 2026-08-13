package wiring

import (
	"context"
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
	"github.com/DrizzDev/platform/internal/integration/domain/agent"
)

// command is the Drizz `/author` slash command, embedded so connect writes the same authoring guide on every machine.
//
//go:embed asset/author.md
var command []byte

// marker identifies a command file Drizz owns. A file that lacks it is the person's own command and is never written
// over or removed.
const marker = "drizz:managed"

// Command installs the Drizz `/author` command for an agent that has a command surface. It reports a conflict, without
// writing, when a command file already exists that Drizz does not own, so the caller can tell the person their own
// command was kept. An agent with no command surface is a no-op.
func (store Store) Command(scope context.Context, job connect.Task) (conflict bool, failure error) {
	commanding := job.Agent.Commanding()
	if !commanding.Supported() {
		return false, nil
	}
	scope, gauge := store.begin(scope, "command")
	defer func() { gauge.close(scope, failure) }()

	path, failure := store.commandpath(commanding)
	if failure != nil {
		return false, failure
	}
	owned, exists, failure := store.owned(path)
	if failure != nil {
		return false, failure
	}
	if exists && !owned {
		return true, nil
	}
	return false, store.write(change{path: path, raw: command})
}

// Uncommand removes the Drizz `/author` command, leaving a command file Drizz does not own in place.
func (store Store) Uncommand(scope context.Context, target agent.Agent) (failure error) {
	commanding := target.Commanding()
	if !commanding.Supported() {
		return nil
	}
	scope, gauge := store.begin(scope, "uncommand")
	defer func() { gauge.close(scope, failure) }()

	path, failure := store.commandpath(commanding)
	if failure != nil {
		return failure
	}
	owned, exists, failure := store.owned(path)
	if failure != nil {
		return failure
	}
	if !exists || !owned {
		return nil
	}
	if failure := os.Remove(path); failure != nil {
		return connect.Locked{}
	}
	return nil
}

// commandpath resolves the agent's command file, refusing a symlinked or irregular node so a planted file cannot
// redirect the write.
func (store Store) commandpath(commanding agent.Commanding) (string, error) {
	root, failure := store.anchor(commanding.Base())
	if failure != nil {
		return "", failure
	}
	path := filepath.Join(append([]string{root}, commanding.Segments()...)...)
	if failure := store.screen(path); failure != nil {
		return "", failure
	}
	return path, nil
}

// owned reports whether the command file exists and carries the Drizz marker. A file without the marker belongs to the
// person and must not be overwritten or removed.
func (Store) owned(path string) (owned bool, exists bool, failure error) {
	raw, failure := os.ReadFile(path)
	if errors.Is(failure, os.ErrNotExist) {
		return false, false, nil
	}
	if failure != nil {
		return false, false, connect.Locked{}
	}
	return strings.Contains(string(raw), marker), true, nil
}
