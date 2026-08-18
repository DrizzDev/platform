// Package binary resolves the path of the running Drizz executable, so the entry the installer writes into an agent's
// configuration launches exactly this installed binary rather than whatever a lookup on the search path might find.
package binary

import "os"

// Resolver reports the running executable's path.
type Resolver struct {
	executable func() (string, error)
}

func New() Resolver {
	return Resolver{executable: os.Executable}
}

// Locate returns the path of the running executable as it was invoked. It deliberately does not resolve symlinks: a
// package manager such as Homebrew installs a stable launcher symlink and a version-pinned target, and the entry the
// installer writes must point at the launcher so it keeps working after an upgrade replaces the target.
func (resolver Resolver) Locate() (string, error) {
	locate := resolver.executable
	if locate == nil {
		locate = os.Executable
	}
	return locate()
}
