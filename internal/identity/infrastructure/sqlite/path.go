package sqlite

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	folder = "drizz"
	name   = "identity.db"
)

type Location struct {
	Root string
}

func (location Location) Resolve() (string, error) {
	root := location.Root
	if root == "" {
		resolved, failure := os.UserConfigDir()
		if failure != nil {
			return "", failure
		}
		root = resolved
	}
	directory := filepath.Join(root, folder)
	if failure := os.MkdirAll(directory, 0o700); failure != nil {
		return "", failure
	}
	if failure := location.screen(directory); failure != nil {
		return "", failure
	}
	target := filepath.Join(directory, name)
	if failure := location.inspect(target); failure != nil {
		return "", failure
	}
	return target, nil
}

// screen rejects a data directory that is a symlink, closing a redirection attack.
func (location Location) screen(path string) error {
	info, failure := os.Lstat(path)
	if failure != nil {
		return failure
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("identity data directory must not be a symlink")
	}
	return nil
}

// inspect rejects a pre-existing database file that is a symlink or is otherwise
// not a regular file, so a planted node cannot redirect writes elsewhere.
func (location Location) inspect(path string) error {
	info, failure := os.Lstat(path)
	switch {
	case errors.Is(failure, os.ErrNotExist):
		return nil
	case failure != nil:
		return failure
	case info.Mode()&os.ModeSymlink != 0:
		return errors.New("identity database file must not be a symlink")
	case !info.Mode().IsRegular():
		return errors.New("identity database file must be a regular file")
	}
	return nil
}
