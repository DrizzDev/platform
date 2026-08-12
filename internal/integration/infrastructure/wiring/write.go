package wiring

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/DrizzDev/platform/internal/integration/application/connect"
)

// suffix is appended to the original file to hold the backup taken before a change, so a person can restore the
// pre-Drizz configuration by hand.
const suffix = ".drizz.bak"

// change is one configuration file to publish: its resolved path and the fully rendered bytes.
type change struct {
	path string
	raw  []byte
}

// stamp backs up the current file, then publishes the new content atomically, so a crash mid-write leaves either the
// original (restorable from the backup) or the fully written file, never a half-written one.
func (store Store) stamp(edit change) error {
	if failure := store.backup(edit.path); failure != nil {
		return failure
	}
	return store.write(edit)
}

func (Store) backup(path string) error {
	raw, failure := os.ReadFile(path)
	if errors.Is(failure, os.ErrNotExist) {
		return nil
	}
	if failure != nil {
		return connect.Locked{}
	}
	if failure := os.WriteFile(path+suffix, raw, 0o600); failure != nil {
		return connect.Locked{}
	}
	return nil
}

func (Store) write(item change) error {
	directory := filepath.Dir(item.path)
	if failure := os.MkdirAll(directory, 0o700); failure != nil {
		return connect.Locked{}
	}
	temporary, failure := os.CreateTemp(directory, "drizz-*")
	if failure != nil {
		return connect.Locked{}
	}
	defer func() { _ = os.Remove(temporary.Name()) }()

	if _, failure := temporary.Write(item.raw); failure != nil {
		_ = temporary.Close()
		return connect.Locked{}
	}
	if failure := temporary.Close(); failure != nil {
		return connect.Locked{}
	}
	return os.Rename(temporary.Name(), item.path)
}
