// Package binary resolves the path of the running Drizz executable, so the entry the installer writes into an agent's
// configuration launches exactly this installed binary rather than whatever a lookup on the search path might find.
package binary

import (
	"os"
	"path/filepath"
)

// Resolver reports the running executable's path.
type Resolver struct{}

func New() Resolver {
	return Resolver{}
}

// Locate returns the absolute, symlink-resolved path of the running executable. If the symlink cannot be resolved the
// raw path is returned, since a launchable path is better than none.
func (Resolver) Locate() (string, error) {
	path, failure := os.Executable()
	if failure != nil {
		return "", failure
	}
	if resolved, broken := filepath.EvalSymlinks(path); broken == nil {
		return resolved, nil
	}
	return path, nil
}
