package carrier

import (
	"errors"
	"os"
)

// publish writes the verified helper to the protected location atomically: it screens the root against a symlink
// redirect, writes to a bounded temporary file in the same directory, marks it executable, and renames it over the
// target so a crash between write and rename leaves only a temporary, never a half-written helper.
func (carrier Carrier) publish(item asset) error {
	if failure := os.MkdirAll(carrier.root, 0o700); failure != nil {
		return failure
	}
	if failure := carrier.screen(); failure != nil {
		return failure
	}
	temporary, failure := os.CreateTemp(carrier.root, executable+"-*")
	if failure != nil {
		return failure
	}
	defer func() { _ = os.Remove(temporary.Name()) }()

	if _, failure := temporary.Write(item.bytes); failure != nil {
		_ = temporary.Close()
		return failure
	}
	if failure := temporary.Chmod(0o700); failure != nil {
		_ = temporary.Close()
		return failure
	}
	if failure := temporary.Close(); failure != nil {
		return failure
	}
	return os.Rename(temporary.Name(), carrier.target())
}

// screen rejects a root that is a symlink or otherwise not a directory, so a planted node cannot redirect the helper
// write elsewhere.
func (carrier Carrier) screen() error {
	info, failure := os.Lstat(carrier.root)
	if failure != nil {
		return failure
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("device helper cache directory must not be a symlink")
	}
	if !info.IsDir() {
		return errors.New("device helper cache path must be a directory")
	}
	return nil
}
