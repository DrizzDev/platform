package bridge

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// verify rejects a tampered or swapped binary before exec: the path must resolve to a file no other user can write,
// and its digest must match the pin.
func (options Options) verify() error {
	resolved, failure := filepath.EvalSymlinks(options.Location)
	if failure != nil {
		return errors.New("device sidecar path could not be resolved")
	}
	info, failure := os.Stat(resolved)
	if failure != nil {
		return errors.New("device sidecar is not present")
	}
	if info.IsDir() || info.Mode().Perm()&0o022 != 0 {
		return errors.New("device sidecar path is not a trusted file")
	}
	digest, failure := options.digest(resolved)
	if failure != nil {
		return failure
	}
	if digest != options.Digest {
		return errors.New("device sidecar failed integrity verification")
	}
	return nil
}

func (options Options) digest(resolved string) (string, error) {
	file, failure := os.Open(resolved)
	if failure != nil {
		return "", errors.New("device sidecar could not be read")
	}

	defer func() { _ = file.Close() }()
	hash := sha256.New()

	if _, failure := io.Copy(hash, file); failure != nil {
		return "", errors.New("device sidecar could not be read")
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
