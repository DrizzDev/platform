package artifact

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/DrizzDev/platform/internal/capture/domain/digest"
)

// stored is one published object as seen on disk: its content address, size, and last write time.
type stored struct {
	size     int64
	modified time.Time
	digest   digest.Digest
}

func (store Store) walk() ([]stored, error) {
	root := filepath.Join(store.root, objects)
	shards, failure := os.ReadDir(root)

	if errors.Is(failure, os.ErrNotExist) {
		return nil, nil
	}
	if failure != nil {
		return nil, failure
	}
	var found []stored
	for _, shard := range shards {
		if !shard.IsDir() {
			continue
		}
		gathered, failure := store.shard(filepath.Join(root, shard.Name()))
		if failure != nil {
			return nil, failure
		}
		found = append(found, gathered...)
	}
	return found, nil
}

// shard reads one content-address prefix directory. A file whose name is not a valid digest is not ours and is skipped.
func (Store) shard(path string) ([]stored, error) {
	files, failure := os.ReadDir(path)

	if failure != nil {
		return nil, failure
	}

	var found []stored
	for _, file := range files {
		key, failure := digest.New(file.Name())
		if failure != nil {
			continue
		}
		info, failure := file.Info()
		if failure != nil {
			continue
		}
		found = append(found, stored{digest: key, modified: info.ModTime(), size: info.Size()})
	}
	return found, nil
}
