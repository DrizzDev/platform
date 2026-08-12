package filesystem

import (
	"os"
	"path/filepath"
)

const (
	folder   = "drizz"
	captures = "captures"
)

type Scratch struct{}

func New() Scratch {
	return Scratch{}
}

type File struct {
	Extension string
	Content   []byte
}

// Save writes a captured file to Drizz's per-user cache directory and returns its path, giving an agent or a person a
// stable, predictable location to work with rather than the opaque, auto-cleaned operating-system temporary directory.
func (scratch Scratch) Save(file File) (string, error) {
	directory, failure := scratch.directory()
	if failure != nil {
		return "", failure
	}
	handle, failure := os.CreateTemp(directory, "capture-*."+file.Extension)
	if failure != nil {
		return "", failure
	}
	_, failure = handle.Write(file.Content)
	closing := handle.Close()
	if failure != nil {
		return "", failure
	}
	return handle.Name(), closing
}

// directory is Drizz's per-user captures directory, created if missing.
func (Scratch) directory() (string, error) {
	cache, failure := os.UserCacheDir()
	if failure != nil {
		return "", failure
	}
	path := filepath.Join(cache, folder, captures)
	if failure := os.MkdirAll(path, 0o700); failure != nil {
		return "", failure
	}
	return path, nil
}
