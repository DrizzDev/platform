package artifact

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

const (
	stage   = "tmp"
	objects = "objects"
)

// locate is the content-addressed path for a digest, sharded by its first two characters to keep any directory small.
func (store Store) locate(key digest.Digest) string {
	value := key.String()
	return filepath.Join(store.root, objects, value[:2], value)
}

// prepare creates the object and staging directories and rejects a root that is a symlink,
// closing a redirection attack; it returns the staging directory.
func (options Options) prepare() (string, error) {
	if failure := os.MkdirAll(filepath.Join(options.Root, objects), 0o700); failure != nil {
		return "", failure
	}

	temp := filepath.Join(options.Root, stage)
	if failure := os.MkdirAll(temp, 0o700); failure != nil {
		return "", failure
	}

	info, failure := os.Lstat(options.Root)
	if failure != nil {
		return "", failure
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return "", errors.New("artifact root must not be a symlink")
	}

	return temp, nil
}
