// Package carrier turns the device helper carried inside the platform binary into a verified, executable file on
// disk. Installing the platform binary installs the helper; the first device call materializes it here, so a person
// never installs, locates, or configures the helper separately.
package carrier

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// folder is the per-user cache subdirectory that holds Drizz's extracted, re-derivable artifacts.
	folder = "drizz"
	// executable is the base name of the extracted device helper.
	executable = "bridge"
	// placeholder is the leading marker of the committed development asset. A release build replaces the asset with
	// the compiled helper; until then the carrier reports no helper is carried so the device path stays cleanly
	// unprepared rather than running a stub.
	placeholder = "drizz-device-helper-placeholder"
)

// Carrier extracts the embedded device helper to a protected per-user location and reuses it. The extracted file is a
// cache: it is re-derivable from the binary, so it lives under the user cache directory rather than durable storage.
type Carrier struct {
	root string
}

func New() (Carrier, error) {
	cache, failure := os.UserCacheDir()
	if failure != nil {
		return Carrier{}, failure
	}
	return Carrier{root: filepath.Join(cache, folder)}, nil
}

// Materialize verifies the carried helper against its pinned digest, writes it once to the protected location, and
// returns the path to run. A verified copy already on disk is reused; a digest mismatch is refused, never run. A build
// that still carries the placeholder reports absent, so the device path stays cleanly unprepared.
func (carrier Carrier) Materialize() (string, error) {
	item := asset{bytes: bridge, digest: fingerprint}
	if bytes.HasPrefix(item.bytes, []byte(placeholder)) {
		return "", errors.New("no device helper is embedded in this build")
	}
	return carrier.deliver(item)
}

// Digest is the pinned hex digest of the carried helper, so a caller can verify it again at the point of use.
func (carrier Carrier) Digest() string {
	return strings.TrimSpace(fingerprint)
}

func (carrier Carrier) deliver(item asset) (string, error) {
	if !item.verify(item.bytes) {
		return "", errors.New("embedded device helper failed its integrity check")
	}
	if carrier.ready(item) {
		return carrier.target(), nil
	}
	if failure := carrier.publish(item); failure != nil {
		return "", failure
	}
	return carrier.target(), nil
}

// target is the path of the extracted helper. Windows requires the executable suffix.
func (carrier Carrier) target() string {
	name := executable
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(carrier.root, name)
}

// ready reports whether a verified copy of the helper is already extracted, so a matching copy is reused rather than
// rewritten on every start.
func (carrier Carrier) ready(item asset) bool {
	present, failure := os.ReadFile(carrier.target())
	if failure != nil {
		return false
	}
	return item.verify(present)
}
