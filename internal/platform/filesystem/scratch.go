package filesystem

import "os"

// Scratch writes short-lived files for a person to open, such as a captured screenshot. It is the approved boundary
// for this kind of local file output, so transports and application code do not touch the filesystem directly.
type Scratch struct{}

func New() Scratch {
	return Scratch{}
}

// File is content to write for a person to open: its bytes and the extension its temporary file should carry.
type File struct {
	Extension string
	Content   []byte
}

// Save writes the file to a new temporary path and returns that path.
func (Scratch) Save(file File) (string, error) {
	handle, failure := os.CreateTemp("", "drizz-*."+file.Extension)
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
